package async

import (
	"consumer-service/internal/kafka"
	"encoding/json"
	kafkaConfig "github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/configs"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/kafka/model"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/redis"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"

	"github.com/spf13/viper"
)

var (
	redisClient redis.Redis
	bufferSize  int
)

const (
	orderSuccess = "success"
	orderFail    = "fail"
)

type Consumer struct {
}

func init() {
	redisClient.Initialize()

	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$PROJECT_PATH/pkg/configs/")

	var configuration kafkaConfig.Configuration

	if err := viper.ReadInConfig(); err != nil {
		logger.KafkaConsumer.Fatalf("Error when reading kafka config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.KafkaConsumer.Fatalf("Unable to decode the kafka config into struct, %v", err)
	}

	bufferSize = configuration.Kafka.BufferSize
}

func NewAsyncConsumer() kafka.BuyEventConsumer {
	return &Consumer{}
}

func (_ *Consumer) StartConsume(brokerList []string, topics []string, group string, config *sarama.Config) {
	config.ClientID = group
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		logger.KafkaConsumer.Fatalf("Failed to start consumer: %v\n", err)
	}

	var (
		messages = make(chan *sarama.ConsumerMessage, bufferSize)
		closing  = make(chan struct{})
		wg       sync.WaitGroup
	)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, os.Interrupt)
		<-signals
		logger.KafkaConsumer.Println("Initiating shutdown of consumer...")
		close(closing)
	}()

	partitions, err := consumer.Partitions(topics[0])
	if err != nil {
		logger.KafkaConsumer.Fatalf("Failed to get partition of topic %v: %v\n", topics[0], err)
	}
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topics[0], partition, sarama.OffsetNewest)
		if err != nil {
			logger.KafkaConsumer.Fatalf("Failed to start consumer for partition %d: %s\n", partition, err)
		}

		go func(pc sarama.PartitionConsumer) {
			<-closing
			pc.AsyncClose()
		}(pc)

		logger.KafkaConsumer.Printf("Consumer connected to the partition %v\n", partition)

		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for message := range pc.Messages() {
				messages <- message
			}
		}(pc)
	}

	go func() {
		for msg := range messages {
			consumeMessage(msg)
		}
	}()

	wg.Wait()
	logger.KafkaConsumer.Printf("Done consuming topic %v\b", topics[0])
	close(messages)

	if err := consumer.Close(); err != nil {
		logger.KafkaConsumer.Printf("Failed to close consumer: %v\n", err)
	}
}

func consumeMessage(message *sarama.ConsumerMessage) {
	logger.KafkaConsumer.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, partition = %v\n",
		string(message.Value), message.Timestamp, message.Topic, message.Partition)

	purchaseMsg := &model.PurchaseMessage{}
	err := json.Unmarshal(message.Value, purchaseMsg)
	if err != nil {
		logger.KafkaConsumer.Printf("Convert purchase message struct to bytes occurs error: %v\n", err)
		return
	}

	// query `Cart` table and update the `Product` table
	carts, err := mariadb.FindAllCartsByCardIDs(purchaseMsg.CartIDs)
	if err != nil {
		logger.KafkaConsumer.Printf("Get carts in consumer occurs error: %v\n", err)
		return
	}

	newOrder := mariadb.Order{}
	newOrder.Initialize(purchaseMsg.AccountID, 0)
	if _, err = newOrder.SaveOrder(); err != nil {
		logger.KafkaConsumer.Printf("Create order occurs error: %v\n", err)
		return
	}

	status, err := updateTables(carts, newOrder, purchaseMsg.AccountID)
	if err != nil {
		logger.KafkaConsumer.Printf("Update purchase related tables occurs error: %v\n", err)
		return
	}

	if err = publishToRedis(purchaseMsg.RedisChannel, status); err != nil {
		logger.KafkaConsumer.Printf("Publish the results to Redis occurs error: %v\n", err)
		return
	}

	go func() {
		if status == orderSuccess {
			if products, err := mariadb.FindAllProducts(); err != nil {
				logger.KafkaConsumer.Warnf("Get all products occurs error: %v\n", err)
			} else if err = redisClient.SetAllProducts(products); err != nil {
				logger.KafkaConsumer.Warnf("Cache all products occurs error: %v\n", err)
			}
		}

		if orders, err := mariadb.FindAllOrdersByAccountID(purchaseMsg.AccountID); err != nil {
			logger.KafkaConsumer.Warnf("Get all orders by accountID occurs error: %v\n", err)
		} else if err = redisClient.SetAllOrdersByAccountID(purchaseMsg.AccountID, orders); err != nil {
			logger.KafkaConsumer.Warnf("Cache all orders by accountID occurs error: %v\n", err)
		}
	}()
}

func updateTables(carts []*mariadb.Cart, newOrder mariadb.Order, accountID int) (string, error) {
	account, err := mariadb.FindAccountByID(accountID)
	if err != nil {
		logger.KafkaConsumer.Warnf("Get account by id occurs error: %v\n", err)
	}

	var amounts = 0
	for _, cart := range carts {
		if amounts > account.Amount {
			break
		}

		product, err := mariadb.PurchaseProduct(cart.ProductID, cart.Quantity)
		if err != nil {
			logger.KafkaConsumer.Printf("Purchase occurs error: %v", err)
			continue
		}

		orderItem := mariadb.OrderItem{}
		orderItem.Initialize(newOrder.ID, cart.ProductID, cart.Quantity)
		if _, err = orderItem.SaveOrderItem(); err != nil {
			logger.KafkaConsumer.Printf("Save the orderItem with product %v occurs error: %v", product.Name, err)
		}

		amounts += product.Price * cart.Quantity
	}

	if amounts == 0 || amounts > account.Amount {
		newOrder.Status = orderFail
	} else {
		newOrder.Amount = amounts

		account.Amount -= amounts
		if _, err = account.UpdateAccount(accountID); err != nil {
			logger.KafkaConsumer.Warnf("Update account %v occurs error: %v\n", accountID, err)
		}
	}

	if _, err := newOrder.UpdateOrder(); err != nil {
		logger.KafkaConsumer.Warnf("Update order %v occurs error: %v\n", newOrder.ID, err)
		return "", err
	}

	return newOrder.Status, nil
}

func publishToRedis(channel string, message string) error {
	if err := redisClient.Publish(channel, message); err != nil {
		return err
	}

	return nil
}

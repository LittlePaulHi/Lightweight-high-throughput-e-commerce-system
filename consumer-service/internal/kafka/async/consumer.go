package async

import (
	"consumer-service/internal/kafka"
	"context"
	"encoding/json"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/kafka/model"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/redis"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

var (
	redisClient redis.Redis
)

const (
	orderSuccess = "success"
	orderFail    = "fail"
)

type Consumer struct {
	Ready chan bool
}

func init() {
	redisClient.Initialize()
}

func NewAsyncConsumer() kafka.BuyEventConsumer {
	return &Consumer{
		Ready: make(chan bool),
	}
}

func (consumer *Consumer) StartConsume(brokerList []string, topics []string, group string, config *sarama.Config) {
	config.ClientID = group
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokerList, group, config)
	if err != nil {
		logger.KafkaConsumer.Panicf("Error when creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, consumer); err != nil {
				logger.KafkaConsumer.Panicf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}

			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	logger.KafkaConsumer.Println("Sarama consumer start running...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logger.KafkaConsumer.Println("terminating: context cancelled")
	case <-sigterm:
		logger.KafkaConsumer.Println("terminating: via signal")
	}

	cancel()
	wg.Wait()

	if err = client.Close(); err != nil {
		logger.KafkaConsumer.Panicf("Error occurs when closing client: %v", err)
	}
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		logger.KafkaConsumer.Printf("Message claimed: value = %s, timestamp = %v, topic = %s",
			string(message.Value), message.Timestamp, message.Topic)

		purchaseMsg := &model.PurchaseMessage{}
		err := json.Unmarshal(message.Value, purchaseMsg)
		if err != nil {
			logger.KafkaConsumer.Printf("Convert purchase message struct to bytes occurs error: %v\n", err)
			continue
		}

		// query `Cart` table and update the `Product` table
		carts, err := mariadb.FindAllCartsByCardIDs(purchaseMsg.CartIDs)
		if err != nil {
			logger.KafkaConsumer.Printf("Get carts in consumer occurs error: %v\n", err)
			continue
		}

		newOrder := mariadb.Order{}
		newOrder.Initialize(purchaseMsg.AccountID, 0)
		if _, err = newOrder.SaveOrder(); err != nil {
			logger.KafkaConsumer.Printf("Create order occurs error: %v\n", err)
			continue
		}

		status, err := updateTables(carts, newOrder, purchaseMsg.AccountID)
		if err != nil {
			logger.KafkaConsumer.Printf("Update purchase related tables occurs error: %v\n", err)
			continue
		}

		if err = publishToRedis(purchaseMsg.RedisChannel, status); err != nil {
			logger.KafkaConsumer.Printf("Publish the results to Redis occurs error: %v\n", err)
			continue
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

		session.MarkMessage(message, "")
	}

	return nil
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

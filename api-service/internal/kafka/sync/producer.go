package sync

import (
	"encoding/json"
	config "github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/configs"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/kafka/model"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/redis"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

var (
	redisClient    redis.Redis
	topics         []string
	brokerList     []string
	flushFrequency time.Duration
)

func init() {
	redisClient.Initialize()

	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$PROJECT_PATH/pkg/configs/")

	var configuration config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		logger.KafkaProducer.Panicf("Error when reading kafka config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.KafkaProducer.Panicf("Unable to decode the kafka config into struct, %v", err)
	}

	brokerList = configuration.Kafka.BrokerList
	topics = configuration.Kafka.Topics
	flushFrequency = configuration.Kafka.FlushFrequency
}

type Kafka struct {
	Producer sarama.SyncProducer
}

func CrateNewSyncProducer() sarama.SyncProducer {
	// setup log to stdout for debug
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	sconfig := sarama.NewConfig()
	sconfig.Producer.Partitioner = sarama.NewRandomPartitioner
	sconfig.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	sconfig.Producer.Retry.Max = 10
	sconfig.Producer.Return.Successes = true
	sconfig.Producer.Flush.Frequency = flushFrequency * time.Millisecond

	producer, err := sarama.NewSyncProducer(brokerList, sconfig)
	if err != nil {
		logger.KafkaProducer.Warnf("Failed to start Sarama SyncProducer, %v\n", err)
	}

	return producer
}

func (kafka *Kafka) PublishBuyEvent(accountID int, cartIDs []int) (string, error) {
	rand.Seed(time.Now().UnixNano())
	purchaseMsg := model.PurchaseMessage{
		RedisChannel: strconv.Itoa(accountID) + "." + time.Now().String() + "." + strconv.Itoa(rand.Int()),
		AccountID:    accountID,
		CartIDs:      cartIDs,
	}
	purchaseMsgBytes, err := json.Marshal(purchaseMsg)
	if err != nil {
		logger.KafkaProducer.Panicf("Convert purchase message struct to bytes occurs error: %v\n", err)
		return "", err
	}

	// TODO: use Flush to batched up and send the message

	msg := &sarama.ProducerMessage{
		Topic: topics[0],
		Value: sarama.ByteEncoder(purchaseMsgBytes),
	}

	partition, offset, err := kafka.Producer.SendMessage(msg)
	if err != nil {
		logger.KafkaProducer.Panicf("Failed to store your message, %v\n", err)
		return "", err
	}
	logger.KafkaProducer.Printf("Purchase message is stored with unique identifier important partition: %v, offset: %v\n", partition, offset)

	payload, err := redisClient.SubscribeAndReceive(purchaseMsg.RedisChannel)
	if err != nil {
		logger.RedisLog.Warnf("SubscribeAndReceive the result of buy event occurs error, %v", err)
		return "", err
	}

	return payload, nil
}

func (kafka *Kafka) Close() error {
	if err := kafka.Producer.Close(); err != nil {
		logger.KafkaProducer.Warnf("Failed to close the kafka producer", err)
	}

	return nil
}

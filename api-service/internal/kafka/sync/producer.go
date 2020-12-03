package sync

import (
	"encoding/json"
	config "github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/configs"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/kafka/model"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

var brokerList []string

func init() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$PROJECT_PATH/pkg/configs/")

	var configuration config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error when reading kafka config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode the kafka config into struct, %v", err)
	}

	brokerList = configuration.Kafka.BrokerList
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

	producer, err := sarama.NewSyncProducer(brokerList, sconfig)
	if err != nil {
		log.Printf("Failed to start Sarama SyncProducer, %v\n", err)
	}

	return producer
}

func (kafka *Kafka) Publish(topic string, cartIDs []int) error {
	purchaseMsg := model.PurchaseMessage{CartIDs: cartIDs}
	purchaseMsgBytes, err := json.Marshal(purchaseMsg)
	if err != nil {
		log.Panicf("Convert purchase message struct to bytes occurs error: %v\n", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.ByteEncoder(purchaseMsgBytes),
	}

	partition, offset, err := kafka.Producer.SendMessage(msg)
	if err != nil {
		log.Panicf("Failed to store your message, %v\n", err)
		return err
	}
	log.Printf("Purchase message is stored with unique identifier important partition: %v, offset: %v\n", partition, offset)

	return nil
}

func (kafka *Kafka) Close() error {
	if err := kafka.Producer.Close(); err != nil {
		log.Fatal("Failed to close the kafka producer", err)
	}

	return nil
}

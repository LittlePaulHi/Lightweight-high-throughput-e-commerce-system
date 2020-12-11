package main

import (
	syncKafka "consumer-service/internal/kafka/sync"
	kafkaConfig "github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/configs"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"context"
	"log"
	"os"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

var (
	brokerList []string
	topics     []string
	group      string
	assignor   string
	verbose    bool
)

const (
	RoundRobin = "roundrobin"
	Sticky     = "sticky"
	Range      = "range"
)

func init() {
	logger.SetLogLevel(logrus.DebugLevel)

	mariadb.Setup()
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$PROJECT_PATH/pkg/configs/")

	var configuration kafkaConfig.Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error when reading kafka config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode the kafka config into struct, %v", err)
	}

	brokerList = configuration.Kafka.BrokerList
	topics = configuration.Kafka.Topics
	group = configuration.Kafka.Group
	assignor = configuration.Kafka.Assignor
	verbose = configuration.Kafka.Verbose
}

func main() {
	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama consumer] ", log.LstdFlags)
	}

	config := sarama.NewConfig()

	switch assignor {
	case RoundRobin:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case Sticky:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case Range:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	consumer := syncKafka.Consumer{
		Ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokerList, group, config)
	if err != nil {
		log.Panicf("Error when creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}

			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer start running...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error occurs when closing client: %v", err)
	}
}

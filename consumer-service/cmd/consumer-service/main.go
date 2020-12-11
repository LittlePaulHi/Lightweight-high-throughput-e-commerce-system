package main

import (
	syncKafka "consumer-service/internal/kafka/sync"
	kafkaConfig "github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/configs"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"

	"github.com/sirupsen/logrus"

	"log"
	"os"

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
		logger.KafkaConsumer.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	consumer := syncKafka.Consumer{
		Ready: make(chan bool),
	}
	consumer.StartConsume(brokerList, topics, group, config)
}

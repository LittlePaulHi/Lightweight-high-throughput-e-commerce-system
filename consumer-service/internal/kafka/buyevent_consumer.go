package kafka

import (
	"github.com/Shopify/sarama"
)

type BuyEventConsumer interface {
	StartConsume(brokerList []string, topics []string, group string, config *sarama.Config)
}

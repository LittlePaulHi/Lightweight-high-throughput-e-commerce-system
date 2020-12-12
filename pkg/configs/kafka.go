package config

import (
	"time"
)

type KafkaConfiguration struct {
	BrokerList     []string
	Topics         []string
	BufferSize     int
	Group          string
	Assignor       string
	Verbose        bool
	FlushFrequency time.Duration // for producer
	ConsumerType   string
}

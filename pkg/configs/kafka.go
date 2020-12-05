package config

type KafkaConfiguration struct {
	BrokerList []string
	Topics     []string
	Group      string
	Assignor   string
	Verbose    bool
}

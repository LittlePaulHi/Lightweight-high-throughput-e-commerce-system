package config

type Configuration struct {
	Mariadb MariadbConfiguration
	Kafka   KafkaConfiguration
	Redis   RedisConfiguration
}

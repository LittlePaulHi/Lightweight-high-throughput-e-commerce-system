module consumer-service

go 1.15

require (
	github.com/Shopify/sarama v1.27.2
	github.com/go-redis/redis/v8 v8.4.0 // indirect
	github.com/spf13/viper v1.7.1
	github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg v0.0.0-00010101000000-000000000000
)

replace github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg => ../pkg

package redis

import (
	"context"
	config "github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/configs"
	"log"
	"time"

	"github.com/spf13/viper"

	"github.com/go-redis/redis/v8"
)

var (
	address string
	ctx     = context.Background()
)

type Redis struct {
	rdb *redis.Client
}

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

	address = configuration.Redis.Address
}

func (redisClient *Redis) Initialize() {
	redisClient.rdb = redis.NewClient(&redis.Options{
		Addr:         address,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	if err := redisClient.rdb.Ping(ctx).Err(); err != nil {
		log.Panicf("Unable to connect to redis: %v", err)
	}
}

func (redisClient *Redis) Publish(channel string, message string) error {
	pub := redisClient.rdb.Subscribe(ctx, channel)
	defer func() {
		if err := pub.Close(); err != nil {
			log.Printf("Close publish from Redis occurs error: %v\n", err)
		}
	}()

	// Wait for confirmation that subscription is created before publishing anything.
	if _, err := pub.Receive(ctx); err != nil {
		panic(err)
	}

	if err := redisClient.rdb.Publish(ctx, channel, message).Err(); err != nil {
		panic(err)
	}

	return nil
}

func (redisClient *Redis) SubscribeAndReceive(channel string) (string, error) {
	sub := redisClient.rdb.Subscribe(ctx, channel)
	defer func() {
		if err := sub.Close(); err != nil {
			log.Printf("Close subscribe from Redis occurs error: %v\n", err)
		}
	}()

	// Wait for confirmation that subscription is created before publishing anything.
	if _, err := sub.Receive(ctx); err != nil {
		panic(err)
	}

	msgI, err := sub.ReceiveTimeout(ctx, 5*time.Second)
	if err != nil {
		log.Printf("Subscribe from Redis and reveive message occurs error: %v", err)
		return "", err
	}

	switch msg := msgI.(type) {
	case *redis.Message:
		log.Printf("Redis receive %v from %v", msg.Payload, msg.Channel)
		return msg.Payload, nil
	default:
		panic("Unreached redis message")
	}
}

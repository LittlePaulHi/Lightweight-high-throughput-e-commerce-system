package redis

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

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

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
	address         string
	dataBase        int
	dialTimeout     time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration
	poolSize        int
	poolTimeout     time.Duration
	cacheExpireTime time.Duration

	ctx = context.Background()
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
	dataBase = configuration.Redis.DataBase
	dialTimeout = configuration.Redis.DialTimeout
	readTimeout = configuration.Redis.ReadTimeout
	writeTimeout = configuration.Redis.WriteTimeout
	poolSize = configuration.Redis.PoolSize
	poolTimeout = configuration.Redis.PoolTimeout

	cacheExpireTime = configuration.Redis.CacheExpireTime
}

func (redisClient *Redis) Initialize() {
	redisClient.rdb = redis.NewClient(&redis.Options{
		Addr:         address,
		DialTimeout:  dialTimeout * time.Second,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
		PoolSize:     poolSize,
		PoolTimeout:  poolTimeout * time.Second,
	})

	if err := redisClient.rdb.Ping(ctx).Err(); err != nil {
		log.Panicf("Unable to connect to redis: %v", err)
	}
}

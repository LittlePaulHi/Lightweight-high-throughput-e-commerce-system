package cache

import (
	"github.com/go-redis/redis/v8"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"time"
	"json"
)

type CartCache interface {
	GetAllCartsByAcctID(int) []*mariadb.Cart
	SetAllCartsByAcctID(int, []*mariadb.Cart)
	//GetAllCartsByProdID(int) []*mariadb.Cart
	//SetAllCartsByProdID(int, []*mariadb.Cart)
}

type redisCartCache struct {
	redisHost string
	db        int
	expires   time.Duration
}

func NewRedisCartCache(redisHost string, db int, expires time.Duration) CartCache {
	return &redisCartCache{
		redisHost: redisHost,
		db:        db,
		expires:   expires,
	}
}

func (cache *redisCartCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cache.redisHost,
		Password: "",
		DB: cache.db,
	})
}

func (cache *redisCartCache) GetAllCartsByAcctID(int accID) []*mariadb.Cart {
	client := cache.getClient()

	data, err := client.Get(accID).Result()
	if err != nil {
		return nil
	}

	var carts []*mariadb.Cart

	err = json.Unmarshal([]byte(data), &carts)
	if err != nil {
		//log error here
		return nil
	}

	return carts
}

func (cache *redisCartCache) SetAllCartsByAcctID(int accID, value []*mariadb.Cart)  {
	client := cache.getClient()

	jsonData, err := json.Marshal(value)

	if err != nil {
		// loge error here
		return
	}

	client.Set(accID, jsonData, cache.expires*time.Second)
}

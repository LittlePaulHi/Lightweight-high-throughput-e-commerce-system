package cache

import (
	"context"
	"encoding/json"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"time"

	"github.com/go-redis/redis/v8"
)

type ProductCache interface {
	GetAllProducts() []*mariadb.Product
	SetAllProducts([]*mariadb.Product)
}

type redisProductCache struct {
	redisHost string
	db        int
	expires   time.Duration
}

func NewPostOrderCache(redisHost string, db int, expires time.Duration) ProductCache {
	return &redisProductCache{
		redisHost: redisHost,
		db:        db,
		expires:   expires,
	}
}

func (cache *redisProductCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.redisHost,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisProductCache) GetAllProducts() []*mariadb.Product {
	client := cache.getClient()

	key := "AllProducts"
	data, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}

	var products []*mariadb.Product

	err = json.Unmarshal([]byte(data), &products)
	if err != nil {
		//log error here
		return nil
	}

	return products
}

func (cache *redisProductCache) SetAllProducts(products []*mariadb.Product) {
	client := cache.getClient()

	jsonData, err := json.Marshal(products)

	if err != nil {
		// loge error here
		return
	}

	key := "AllProducts"
	client.Set(context.Background(), key, jsonData, cache.expires*time.Second)
}

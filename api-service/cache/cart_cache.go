package cache

import (
	"github.com/go-redis/redis/v8"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"time"
	"encoding/json"
	"strconv"
	"context"
)

type CartCache interface {
	GetAllCartsByAcctID(accID int) []*mariadb.Cart
	SetAllCartsByAcctID(accID int,carts []*mariadb.Cart)
	GetAllCartsByProdID(int) []*mariadb.Cart
	SetAllCartsByProdID(int, []*mariadb.Cart)
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

func (cache *redisCartCache) GetAllCartsByAcctID(accID int) []*mariadb.Cart {
	client := cache.getClient()


	key := "AllCartsByAcctID" + strconv.Itoa(accID)
	data, err := client.Get(context.Background(), key).Result()
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

func (cache *redisCartCache) GetAllCartsByProdID(prodID int) []*mariadb.Cart {
	client := cache.getClient()


	key := "AllCartsByProdID" + strconv.Itoa(prodID)
	data, err := client.Get(context.Background(), key).Result()
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

func (cache *redisCartCache) SetAllCartsByAcctID(accID int, carts []*mariadb.Cart)  {
	client := cache.getClient()

	jsonData, err := json.Marshal(carts)

	if err != nil {
		// loge error here
		return
	}

	key := "AllCartsByAcctID" + strconv.Itoa(accID)
	client.Set(context.Background(), key, jsonData, cache.expires*time.Second)
}


func (cache *redisCartCache) SetAllCartsByProdID(prodID int, carts []*mariadb.Cart)  {
	client := cache.getClient()

	jsonData, err := json.Marshal(carts)

	if err != nil {
		// loge error here
		return
	}

	key := "AllCartsByProdID" + strconv.Itoa(prodID)
	client.Set(context.Background(), key, jsonData, cache.expires*time.Second)
}
package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"strconv"
	"time"
)

type OrderCache interface {
	GetAllOrdersByAcctID(accID int) []*mariadb.Order
	SetAllOrdersByAcctID(accID int, carts []*mariadb.Order)
	GetAllOrderItemsByOrderID(orderID int) []*mariadb.OrderItem
	SetAllOrderItemsByOrderID(orderID int, items []*mariadb.OrderItem)
}

type redisOrderCache struct {
	redisHost string
	db        int
	expires   time.Duration
}

func NewRedisOrderCache(redisHost string, db int, expires time.Duration) OrderCache {
	return &redisOrderCache{
		redisHost: redisHost,
		db:        db,
		expires:   expires,
	}
}

func (cache *redisOrderCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.redisHost,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisOrderCache) GetAllOrdersByAcctID(accID int) []*mariadb.Order {
	client := cache.getClient()

	key := "AllOrdersByAcctID" + strconv.Itoa(accID)
	data, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}

	var orders []*mariadb.Order

	err = json.Unmarshal([]byte(data), &orders)
	if err != nil {
		//log error here
		return nil
	}

	return orders

}

func (cache *redisOrderCache) SetAllOrdersByAcctID(accID int, orders []*mariadb.Order) {
	client := cache.getClient()

	jsonData, err := json.Marshal(orders)

	if err != nil {
		// loge error here
		return
	}

	key := "AllOrdersByAcctID" + strconv.Itoa(accID)
	client.Set(context.Background(), key, jsonData, cache.expires*time.Second)
}

func (cache *redisOrderCache) GetAllOrderItemsByOrderID(orderID int) []*mariadb.OrderItem {
	client := cache.getClient()

	key := "AllOrderItemsByOrderID" + strconv.Itoa(orderID)
	data, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}

	var orderItems []*mariadb.OrderItem

	err = json.Unmarshal([]byte(data), &orderItems)
	if err != nil {
		//log error here
		return nil
	}

	return orderItems

}

func (cache *redisOrderCache) SetAllOrderItemsByOrderID(orderID int, items []*mariadb.OrderItem) {
	client := cache.getClient()

	jsonData, err := json.Marshal(items)

	if err != nil {
		// loge error here
		return
	}

	key := "AllOrderItemsByOrderID" + strconv.Itoa(orderID)
	client.Set(context.Background(), key, jsonData, cache.expires*time.Second)
}

package redis

import (
	"encoding/json"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"strconv"
	"time"
)

const (
	allProductKey            = "AllProducts"
	allOrderKeyByAccountID   = "AllOrdersByAcctID"
	allOrderItemKeyByOrderID = "AllOrderItemsByOrderID"
)

func (redisClient *Redis) SetAllProducts(products []*mariadb.Product) error {
	jsonData, err := json.Marshal(products)
	if err != nil {
		logger.RedisLog.Warnf("Set cache of allProducts occurs error: %v\n", err)
		return err
	}

	return redisClient.rdb.Set(ctx, allProductKey, jsonData, cacheExpireTime*time.Minute).Err()
}

func (redisClient *Redis) GetAllProducts() []*mariadb.Product {
	data, err := redisClient.rdb.Get(ctx, allProductKey).Result()
	if err != nil {
		return nil
	}

	var products []*mariadb.Product

	err = json.Unmarshal([]byte(data), &products)
	if err != nil {
		logger.RedisLog.Warnf("Get cache of allProducts occurs error: %v\n", err)
		return nil
	}

	return products
}

func (redisClient *Redis) SetAllOrdersByAccountID(accountID int, orders []*mariadb.Order) error {
	jsonData, err := json.Marshal(orders)
	if err != nil {
		logger.RedisLog.Warnf("Set cache of all the orders by accountID %v occurs error: %v\n", accountID, err)
		return err
	}

	key := allOrderKeyByAccountID + strconv.Itoa(accountID)
	return 	redisClient.rdb.Set(ctx, key, jsonData, cacheExpireTime*time.Minute).Err()
}

func (redisClient *Redis) GetAllOrdersByAccountID(accountID int) []*mariadb.Order {
	key := allOrderKeyByAccountID + strconv.Itoa(accountID)
	data, err := redisClient.rdb.Get(ctx, key).Result()
	if err != nil {
		logger.RedisLog.Warnf("Get cache of all the orders by accountID %v occurs error: %v\n", accountID, err)
		return nil
	}

	var orders []*mariadb.Order

	err = json.Unmarshal([]byte(data), &orders)
	if err != nil {
		logger.RedisLog.Warnf("Deserialize cache data to Order struct occurs error: %v\n", err)
		return nil
	}

	return orders
}

func (redisClient *Redis) SetAllOrderItemsByOrderID(orderID int, items []*mariadb.OrderItem) error {
	jsonData, err := json.Marshal(items)
	if err != nil {
		logger.RedisLog.Warnf("Set cache of all the order items by orderID %v occurs error: %v\n", orderID, err)
		return err
	}

	key := allOrderItemKeyByOrderID + strconv.Itoa(orderID)
	return redisClient.rdb.Set(ctx, key, jsonData, cacheExpireTime*time.Minute).Err()
}

func (redisClient *Redis) GetAllOrderItemsByOrderID(orderID int) []*mariadb.OrderItem {
	key := allOrderItemKeyByOrderID + strconv.Itoa(orderID)
	data, err := redisClient.rdb.Get(ctx, key).Result()
	if err != nil {
		logger.RedisLog.Warnf("Get cache of all the order items by orderID %v occurs error: %v\n", orderID, err)
		return nil
	}

	var orderItems []*mariadb.OrderItem

	err = json.Unmarshal([]byte(data), &orderItems)
	if err != nil {
		logger.RedisLog.Warnf("Deserialize cache data to OrderItems struct occurs error: %v\n", err)
		return nil
	}

	return orderItems
}

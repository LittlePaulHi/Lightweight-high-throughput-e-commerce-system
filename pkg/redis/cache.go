package redis

import (
	"encoding/json"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"time"
)

const (
	allProductKey = "AllProducts"
)

func (redisClient *Redis) SetAllProducts(products []*mariadb.Product) error {
	jsonData, err := json.Marshal(products)
	if err != nil {
		logger.RedisLog.Warnf("Set cache of allProducts occurs error: %v\n", err)
		return err
	}

	redisClient.rdb.Set(ctx, allProductKey, jsonData, cacheExpireTime*time.Minute)

	return nil
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

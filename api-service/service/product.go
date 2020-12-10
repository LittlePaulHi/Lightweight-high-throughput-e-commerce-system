package service

import (
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/redis"
)

var (
	productRedisClient redis.Redis
)

func init() {
	productRedisClient.Initialize()
}

// GetProduct from mariadb/redis
func GetProduct(id int) (*mariadb.Product, error) {

	// TODO: cache service

	product, err := mariadb.FindProductByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetAllProducts from mariadb/redis
func GetAllProducts() ([]*mariadb.Product, error) {
	products := productRedisClient.GetAllProducts()
	if products == nil {
		var err error
		products, err = mariadb.FindAllProducts()
		if err != nil {
			logger.APILog.Warn(err)
			return nil, err
		}

		if err = productRedisClient.SetAllProducts(products); err != nil {
			logger.APILog.Warn(err)
		}
	}

	return products, nil
}

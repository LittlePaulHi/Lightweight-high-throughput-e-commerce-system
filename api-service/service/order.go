package service

import (
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/redis"
)

var (
	orderRedisClient redis.Redis
)

func init() {
	orderRedisClient.Initialize()
}

// GetAllOrdersByAccountID gets all orders by the specified account id
func GetAllOrdersByAccountID(accountID int) ([]*mariadb.Order, error) {
	orders := orderRedisClient.GetAllOrdersByAccountID(accountID)
	if orders == nil {
		var err error
		orders, err = mariadb.FindAllOrdersByAccountID(accountID)
		if err != nil {
			return nil, err
		}

		if err = orderRedisClient.SetAllOrdersByAccountID(accountID, orders); err != nil {
			logger.APILog.Warnln(err)
		}
	}

	return orders, nil
}

// GetAllOrderItemsByOrderID gets all orders by the specified order id
func GetAllOrderItemsByOrderID(orderID int) ([]*mariadb.OrderItem, error) {
	orderItems := orderRedisClient.GetAllOrderItemsByOrderID(orderID)
	if orderItems == nil {
		var err error
		orderItems, err = mariadb.FindAllOrderItemsByOrderID(orderID)
		if err != nil {
			return nil, err
		}

		if err = orderRedisClient.SetAllOrderItemsByOrderID(orderID, orderItems); err != nil {
			logger.APILog.Warnln(err)
		}
	}

	return orderItems, nil
}

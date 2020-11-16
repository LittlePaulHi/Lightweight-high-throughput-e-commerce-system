package service

import (
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
)

// GetAllOrdersByAccountID gets all orders by the specified account id
func GetAllOrdersByAccountID(accountID int) ([]*mariadb.Order, error) {

	// TODO: cache service

	orders, err := mariadb.FindAllOrdersByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// GetAllOrderItemsByOrderID gets all orders by the specified order id
func GetAllOrderItemsByOrderID(orderID int) ([]*mariadb.OrderItem, error) {

	// TODO: cache service

	orderItems, err := mariadb.FindAllOrderItemsByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	return orderItems, nil
}

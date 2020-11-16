package service

import (
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
)

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

	// TODO: cache service

	products, err := mariadb.FindAllProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}

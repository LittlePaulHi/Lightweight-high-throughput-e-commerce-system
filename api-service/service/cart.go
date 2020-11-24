package service

import (
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
)



// GetAllCartsByAccountID gets all carts by specified account id
func GetAllCartsByAccountID(accountID int) ([]*mariadb.Cart, error) {

	// TODO: cache service

	carts, err := mariadb.FindAllCartsByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	return carts, nil
}

// GetAllCartsByProductID gets all carts by specified product id
func GetAllCartsByProductID(productID int) ([]*mariadb.Cart, error) {

	// TODO: cache service

	carts, err := mariadb.FindAllCartsByProductID(productID)
	if err != nil {
		return nil, err
	}

	return carts, nil
}

// AllCartsByProductID

// AddCart adds the new cart item (each product map to a unique cart)
func AddCart(cart *mariadb.Cart) (*mariadb.Cart, error) {

	// TODO: cache service

	cart, err := cart.SaveCart()
	if err != nil {
		return nil, err
	}

	return cart, nil
}

// EditCart edits the specified cart
func EditCart(cartID int, accountID int, productID int, quantity int) (*mariadb.Cart, error) {

	// TODO: cache service

	cart, err := mariadb.UpdateCart(cartID, productID, quantity)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

package mariadb

import (
	"gorm.io/gorm"
	"time"
)

// Cart struct used by Mariadb model
type Cart struct {
	ID        int       `gorm:"primary_key;auto_increment;uniqueIndex"`
	AccountID int       `gorm:"not null;uniqueIndex"` // foreign key of Account
	ProductID int       `gorm:"not null;uniqueIndex"` // foreign key of Product
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt
}

// Initialize the cart
func (cart *Cart) Initialize(accountID int, productID int, quantity int) {
	cart.ID = 0
	cart.AccountID = accountID
	cart.ProductID = productID
	cart.Quantity = quantity
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()
}

// SaveCart saves the specified cart
func (cart *Cart) SaveCart() (*Cart, error) {
	err := db.Create(&cart).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return cart, nil
}

// UpdateCart updates the specified cart
func UpdateCart(cartID int, productID int, quantity int) (*Cart, error) {
	db = db.Model(&Order{}).Where("ID = ?", cartID).Take(&Order{}).UpdateColumns(
		map[string]interface{}{
			"ProductID": productID,
			"Quantity":  quantity,
			"UpdatedAt": time.Now(),
		},
	)
	if db.Error != nil {
		return &Cart{}, db.Error
	}

	// check the updated cart
	var cart *Cart
	err := db.Model(&Order{}).Where("ID = ?", cartID).Take(&cart).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return cart, nil
}

// FindAllCartsByAccountID finds all the carts by specified account-id
func FindAllCartsByAccountID(accountID int) ([]*Cart, error) {
	cartItems := []*Cart{}
	err := db.Model(&Cart{}).Where("AccountID = ?", accountID).Find(&cartItems).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return cartItems, nil
}

// FindAllCartsByProductID finds all the carts by specified productID
func FindAllCartsByProductID(productID int) ([]*Cart, error) {
	cartItems := []*Cart{}
	err := db.Model(&Cart{}).Where("ProductID = ?", productID).Find(&cartItems).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return cartItems, nil
}

// DeleteCartByID - soft delete
func DeleteCartByID(id int) (int64, error) {
	db = db.Model(&Cart{}).Where("ID = ?", id).Delete(&Cart{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

// DeleteCartByAccountID - soft delete
func DeleteCartByAccountID(accountID int) (int64, error) {
	db = db.Model(&Cart{}).Where("AccountID = ?", accountID).Delete(&Cart{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

// DeleteCartByProductID - soft delete
func DeleteCartByProductID(productID int) (int64, error) {
	db = db.Model(&Cart{}).Where("ProductID = ?", productID).Delete(&Cart{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

package mariadb

import (
	"time"

	"gorm.io/gorm"
)

// Cart struct used by Mariadb model
type Cart struct {
	ID        int       `gorm:"primaryKey;autoIncrement;uniqueIndex"`
	AccountID int       `gorm:"not null;index"` // foreign key of Account
	ProductID int       `gorm:"not null;index"` // foreign key of Product
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
	tx := db.Model(&Order{}).Where("ID = ?", cartID).Take(&Order{}).UpdateColumns(
		map[string]interface{}{
			"ProductID": productID,
			"Quantity":  quantity,
			"UpdatedAt": time.Now(),
		},
	)
	if tx.Error != nil {
		return &Cart{}, tx.Error
	}

	// check the updated cart
	var cart *Cart
	err := tx.Model(&Order{}).Where("ID = ?", cartID).Take(&cart).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return cart, nil
}

// FindAllCartsByCardIDs find all the carts by specified CartIDs
func FindAllCartsByCardIDs(cartIDs []int) ([]*Cart, error) {
	var cartItems []*Cart
	err := db.Model(&Cart{}).Where("ID IN ?", cartIDs).Find(&cartItems).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return cartItems, nil
}

// FindAllCartsByAccountID finds all the carts by specified account-id
func FindAllCartsByAccountID(accountID int) ([]*Cart, error) {
	cartItems := []*Cart{}
	err := db.Model(&Cart{}).Where("account_id = ?", accountID).Find(&cartItems).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return cartItems, nil
}

// FindAllCartsByProductID finds all the carts by specified productID
func FindAllCartsByProductID(productID int) ([]*Cart, error) {
	cartItems := []*Cart{}
	err := db.Model(&Cart{}).Where("product_id = ?", productID).Find(&cartItems).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return cartItems, nil
}

// DeleteCartByID - soft delete
func DeleteCartByID(id int) (int64, error) {
	tx := db.Model(&Cart{}).Where("ID = ?", id).Delete(&Cart{})
	if tx.Error != nil {
		return 0, tx.Error
	}

	return tx.RowsAffected, nil
}

// DeleteCartByAccountID - soft delete
func DeleteCartByAccountID(accountID int) (int64, error) {
	tx := db.Model(&Cart{}).Where("account_id = ?", accountID).Delete(&Cart{})
	if tx.Error != nil {
		return 0, tx.Error
	}

	return tx.RowsAffected, nil
}

// DeleteCartByProductID - soft delete
func DeleteCartByProductID(productID int) (int64, error) {
	tx := db.Model(&Cart{}).Where("product_id = ?", productID).Delete(&Cart{})
	if tx.Error != nil {
		return 0, tx.Error
	}

	return tx.RowsAffected, nil
}

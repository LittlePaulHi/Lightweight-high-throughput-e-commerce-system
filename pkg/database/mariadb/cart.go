package mariadb

import (
	"gorm.io/gorm"
	"time"
)

type Cart struct {
	ID        int       `gorm:"primary_key;auto_increment;uniqueIndex"`
	AccountID int       `gorm:"not null;uniqueIndex"` // foreign key of Account
	ProductID int       `gorm:"not null;uniqueIndex"` // foreign key of Product
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt
}

func (cart *Cart) Initialize(name string, quantity int) {
	cart.ID = 0
	cart.Quantity = quantity
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()
}

func (cart *Cart) SaveCart(db *gorm.DB) (*Cart, error) {
	err := db.Create(&cart).Error
	if err != nil {
		return &Cart{}, err
	}
	return cart, nil
}

func (cart *Cart) UpdateCart(db *gorm.DB) (*Cart, error) {
	db = db.Model(&Order{}).Where("ID = ?", cart.ID).Take(&Order{}).UpdateColumns(
		map[string]interface{}{
			"ProductID": cart.ProductID,
			"Quantity":  cart.Quantity,
			"UpdatedAt": time.Now(),
		},
	)
	if db.Error != nil {
		return &Cart{}, db.Error
	}

	// check the updated cart
	err := db.Model(&Order{}).Where("ID = ?", cart.ID).Take(&cart).Error
	if err != nil {
		return &Cart{}, err
	}

	return cart, nil
}

func (cart *Cart) FindAllCartsByAccountID(db *gorm.DB, account_id int) (*[]Cart, error) {
	cart_items := []Cart{}
	err := db.Model(&Cart{}).Where("AccountID = ?", account_id).Find(&cart_items).Error
	if err != nil {
		return &[]Cart{}, err
	}

	return &cart_items, nil
}

func (cart *Cart) FindAllCartsByProductID(db *gorm.DB, product_id int) (*[]Cart, error) {
	cart_items := []Cart{}
	err := db.Model(&Cart{}).Where("ProductID = ?", product_id).Find(&cart_items).Error
	if err != nil {
		return &[]Cart{}, err
	}

	return &cart_items, nil
}

/* soft delete */
func (cart *Cart) DeleteCartByID(db *gorm.DB, id int) (int64, error) {
	db = db.Model(&Cart{}).Where("ID = ?", id).Delete(&Cart{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

/* soft delete */
func (cart *Cart) DeleteCartByAccountID(db *gorm.DB, account_id int) (int64, error) {
	db = db.Model(&Cart{}).Where("AccountID = ?", account_id).Delete(&Cart{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

/* soft delete */
func (cart *Cart) DeleteCartByProductID(db *gorm.DB, product_id int) (int64, error) {
	db = db.Model(&Cart{}).Where("ProductID = ?", product_id).Delete(&Cart{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

package mariadb

import (
	"gorm.io/gorm"
	"time"
)

// Product struct used by Mariadb model
type Product struct {
	ID        int       `gorm:"primary_key;auto_increment;uniqueIndex"`
	Name      string    `gorm:"size:100;not null;unique"`
	Price     int       `gorm:"not null"`
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt
}

// Initialize the product
func (product *Product) Initialize(name string, price int, quantity int) {
	product.ID = 0
	product.Name = name
	product.Price = price
	product.Quantity = quantity
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
}

// SaveProduct saves the product
func (product *Product) SaveProduct() (*Product, error) {
	err := db.Create(&product).Error
	if err != nil {
		return &Product{}, err
	}
	return product, nil
}

// UpdateProduct updates product by the specified ID
func (product *Product) UpdateProduct() (*Product, error) {
	db = db.Model(&Product{}).Where("ID = ?", product.ID).Take(&Product{}).UpdateColumns(
		map[string]interface{}{
			"Quantity":  product.Quantity,
			"UpdatedAt": time.Now(),
		},
	)
	if db.Error != nil {
		return &Product{}, db.Error
	}

	// check the updated product
	err := db.Model(&Product{}).Where("ID = ?", product.ID).Take(&product).Error
	if err != nil {
		return &Product{}, err
	}

	return product, nil
}

// FindAllProducts finds all the products
func FindAllProducts() ([]*Product, error) {
	products := []*Product{}
	err := db.Model(&Product{}).Find(&products).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return products, nil
}

// FindProductByID finds a product by the specified ID
func FindProductByID(id int) (*Product, error) {
	var product Product
	err := db.Model(Product{}).Where("ID = ?", id).Take(&product).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &product, nil
}

// DeleteProductByID - soft delete
func DeleteProductByID(id int) (int64, error) {
	db = db.Model(&Product{}).Where("ID = ?", id).Delete(&Product{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

package mariadb

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Product struct used by Mariadb model
type Product struct {
	ID        int       `gorm:"primaryKey;autoIncrement;uniqueIndex"`
	Name      string    `gorm:"size:100;not null;unique"`
	Price     int       `gorm:"not null"`
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt
	Cart      Cart      `gorm:"ForeignKey:ProductID"`
	OrderItem OrderItem `gorm:"ForeignKey:ProductID"`
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
	tx := db.Model(&Product{}).Where("ID = ?", product.ID).Take(&Product{}).UpdateColumns(
		map[string]interface{}{
			"Quantity":  product.Quantity,
			"UpdatedAt": time.Now(),
		},
	)
	if tx.Error != nil {
		return &Product{}, tx.Error
	}

	// check the updated product
	err := tx.Model(&Product{}).Where("ID = ?", product.ID).Take(&product).Error
	if err != nil {
		return &Product{}, err
	}

	return product, nil
}

func PurchaseProduct(productID int, quantity int) (*Product, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	var product Product
	if err := tx.Model(&Product{}).Where("ID = ?", productID).Take(&product).Error; err != nil {
		return nil, err
	}

	if product.Quantity < quantity {
		return nil, errors.New("Product " + product.Name + " already sold out")
	}

	if err := tx.Model(&Product{}).Where("ID = ?", product.ID).Take(&Product{}).UpdateColumns(
		map[string]interface{}{
			"Quantity":  product.Quantity - quantity,
			"UpdatedAt": time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &product, tx.Commit().Error
}

type AllProducts struct {
	ID        int
	Name      string
	Price     int
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FindAllProducts finds all the products
func FindAllProducts() ([]*AllProducts, error) {
	var products []*AllProducts
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
	tx := db.Model(&Product{}).Where("ID = ?", id).Delete(&Product{})
	if tx.Error != nil {
		return 0, tx.Error
	}

	return tx.RowsAffected, nil
}

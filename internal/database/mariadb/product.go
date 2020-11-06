package mariadb

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID        int       `gorm:"primary_key;auto_increment;uniqueIndex"`
	Name      string    `gorm:"size:100;not null;unique"`
	Price     int       `gorm:"not null"`
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (product *Product) Initialize(name string, price int, quantity int) {
	product.ID = 0
	product.Name = name
	product.Price = price
	product.Quantity = quantity
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
}

func (product *Product) SaveProduct(db *gorm.DB) (*Product, error) {
	err := db.Create(&product).Error
	if err != nil {
		return &Product{}, err
	}
	return product, nil
}

func (product *Product) UpdateProduct(db *gorm.DB) (*Product, error) {
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

func (product *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
	products := []Product{}
	err := db.Model(&Product{}).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}

	return &products, nil
}

func (product *Product) FindProductByID(db *gorm.DB, id int) (*Product, error) {
	err := db.Model(Product{}).Where("ID = ?", id).Take(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Product{}, errors.New("Product Not Found")
	}
	if err != nil {
		return &Product{}, err
	}

	return product, nil
}

func (product *Product) DeleteProduct(db *gorm.DB, id int) (int64, error) {
	db = db.Model(&Product{}).Where("ID = ?", id).Take(&Product{}).Delete(&Product{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

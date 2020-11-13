package mariadb

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID         int         `gorm:"primary_key;auto_increment;uniqueIndex"`
	AccountID  int         `gorm:"not null;uniqueIndex"` // foreign key of Account
	Amount     int         `gorm:"not null"`
	Status     string      `gorm:"type:enum('success', 'fail', 'return');default:'success'"`
	CreatedAt  time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	OrderItems []OrderItem `gorm:"ForeignKey:OrderID"`
}

func (order *Order) Initialize(name string, amount int) {
	order.ID = 0
	order.Amount = amount
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
}

func (order *Order) SaveOrder(db *gorm.DB) (*Order, error) {
	err := db.Create(&order).Error
	if err != nil {
		return &Order{}, err
	}
	return order, nil
}

func (order *Order) UpdateOrder(db *gorm.DB) (*Order, error) {
	db = db.Model(&Order{}).Where("ID = ?", order.ID).Take(&Order{}).UpdateColumns(
		map[string]interface{}{
			"Amount":    order.Amount,
			"Status":    order.Status,
			"UpdatedAt": time.Now(),
		},
	)
	if db.Error != nil {
		return &Order{}, db.Error
	}

	// check the updated order
	err := db.Model(&Order{}).Where("ID = ?", order.ID).Take(&order).Error
	if err != nil {
		return &Order{}, err
	}

	return order, nil
}

func (order *Order) FindAllOrders(db *gorm.DB) (*[]Order, error) {
	orders := []Order{}
	err := db.Model(&Order{}).Find(&orders).Error
	if err != nil {
		return &[]Order{}, err
	}

	return &orders, nil
}

func (order *Order) FindOrderByID(db *gorm.DB, id int) (*Order, error) {
	err := db.Model(Order{}).Where("ID = ?", id).Take(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Order{}, errors.New("Order Not Found")
	}
	if err != nil {
		return &Order{}, err
	}

	return order, nil
}

func (order *Order) DeleteOrderByID(db *gorm.DB, id int) (int64, error) {
	db = db.Model(&Order{}).Where("ID = ?", id).Take(&Order{}).Delete(&Order{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

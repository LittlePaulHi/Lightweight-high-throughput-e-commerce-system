package mariadb

import (
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

func (order *Order) SaveOrder() (*Order, error) {
	err := db.Create(&order).Error
	if err != nil {
		return &Order{}, err
	}
	return order, nil
}

func (order *Order) UpdateOrder() (*Order, error) {
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

func FindAllOrders() ([]*Order, error) {
	orders := []*Order{}
	err := db.Model(&Order{}).Find(&orders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return orders, nil
}

func FindOrderByID(id int) (*Order, error) {
	var order *Order
	err := db.Model(Order{}).Where("ID = ?", id).Take(&order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return order, nil
}

func FindAllOrdersByAccountID(accountID int) ([]*Order, error) {
	var order []*Order
	err := db.Model(Order{}).Where("AccountID = ?", accountID).Find(&order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return order, nil
}

func (order *Order) DeleteOrderByID(id int) (int64, error) {
	db = db.Model(&Order{}).Where("ID = ?", id).Take(&Order{}).Delete(&Order{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

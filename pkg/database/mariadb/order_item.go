package mariadb

import (
	"gorm.io/gorm"
	"time"
)

type OrderItem struct {
	ID        int       `gorm:"primaryKey;autoIncrement;uniqueIndex"`
	OrderID   int       `gorm:"not null;uniqueIndex"` // foreign key of Order
	ProductID int       `gorm:"not null;uniqueIndex"` // foreign key of Product
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (orderItem *OrderItem) Initialize(quantity int) {
	orderItem.ID = 0
	orderItem.Quantity = quantity
	orderItem.CreatedAt = time.Now()
	orderItem.UpdatedAt = time.Now()
}

func (orderItem *OrderItem) SaveOrderItem() (*OrderItem, error) {
	err := db.Create(&orderItem).Error
	if err != nil {
		return &OrderItem{}, err
	}
	return orderItem, nil
}

func FindAllOrderItemsByOrderID(orderID int) ([]*OrderItem, error) {
	orderItems := []*OrderItem{}
	err := db.Model(&OrderItem{}).Where("OrderID = ?", orderID).Find(&orderItems).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return orderItems, nil
}

func FindAllOrderItemsByProductID(productID int) ([]*OrderItem, error) {
	orderItems := []*OrderItem{}
	err := db.Model(&OrderItem{}).Where("ProductID = ?", productID).Find(&orderItems).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return orderItems, nil
}

func (orderItem *OrderItem) DeleteOrderItemByID(id int) (int64, error) {
	db = db.Model(&OrderItem{}).Where("ID = ?", id).Take(&OrderItem{}).Delete(&OrderItem{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

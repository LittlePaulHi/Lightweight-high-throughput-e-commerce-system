package mariadb

import (
	"gorm.io/gorm"
	"time"
)

type OrderItem struct {
	ID        int       `gorm:"primary_key;auto_increment;uniqueIndex"`
	OrderID   int       `gorm:"not null"` // foreign key of Order
	ProductID int       `gorm:"not null"` // foreign key of Product
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (order_item *OrderItem) Initialize(name string, quantity int) {
	order_item.ID = 0
	order_item.Quantity = quantity
	order_item.CreatedAt = time.Now()
	order_item.UpdatedAt = time.Now()
}

func (order_item *OrderItem) SaveOrderItem(db *gorm.DB) (*OrderItem, error) {
	err := db.Create(&order_item).Error
	if err != nil {
		return &OrderItem{}, err
	}
	return order_item, nil
}

func (order_item *OrderItem) FindAllOrderItemsByID(db *gorm.DB, id int) (*[]OrderItem, error) {
	order_items := []OrderItem{}
	err := db.Model(&OrderItem{}).Where("ID = ?", id).Find(&order_items).Error
	if err != nil {
		return &[]OrderItem{}, err
	}

	return &order_items, nil
}

func (order_item *OrderItem) DeleteOrderItemByID(db *gorm.DB, id int) (int64, error) {
	db = db.Model(&OrderItem{}).Where("ID = ?", id).Take(&OrderItem{}).Delete(&OrderItem{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

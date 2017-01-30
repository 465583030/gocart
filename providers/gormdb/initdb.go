package gormdb

import (
	"github.com/alioygur/gocart/domain"
	"github.com/jinzhu/gorm"
)

// InitDB creates tables
func InitDB(db *gorm.DB) error {
	return db.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(
		&domain.User{},
		&domain.Address{},
		&domain.Product{},
		&domain.Category{},
		&domain.Image{},
		&domain.Order{},
		&domain.OrderProduct{},
		&domain.OrderHistory{},
		&domain.OrderAddress{},
		&domain.OrderStatus{},
		&domain.PaymentMethod{},
	).Error
}

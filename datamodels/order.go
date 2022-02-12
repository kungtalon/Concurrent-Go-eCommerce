package datamodels

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId      uint `sql:"UserId"`
	ProductId   uint `sql:"ProductId"`
	OrderStatus uint `sql:"OrderStatus"`
}

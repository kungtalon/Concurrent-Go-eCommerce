package datamodels

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId      uint `sql:"UserId"`
	ProductId   uint `sql:"ProductId"`
	OrderStatus int  `sql:"OrderStatus"`
}

const (
	OrderWait = iota
	OrderSuccess
	OrderFailed
)

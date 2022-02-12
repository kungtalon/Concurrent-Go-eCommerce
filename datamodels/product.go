package datamodels

import (
	"github.com/kataras/golog"
	"gorm.io/gorm"
	"strconv"
)

type Product struct {
	gorm.Model
	ProductName  string `json:"ProductName" gorm:"ProductName" lightning:"ProductName"`
	ProductNum   int64  `json:"ProductNum" gorm:"ProductNum" lightning:"ProductNum"`
	ProductImage string `json:"ProductImage" gorm:"ProductImage" lightning:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" gorm:"ProductUrl" lightning:"ProductUrl"`
}

func (p *Product) PrintInfo(log *golog.Logger) {
	log.Debug("Get Product: ID:" + strconv.FormatUint(uint64(p.ID), 10) + " Name: " + p.ProductName)
}

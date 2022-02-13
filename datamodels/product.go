package datamodels

import (
	"fmt"
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
	ProductPrice uint64 `json:"ProductPrice" lightning:"ProductPrice"`
	DealPrice    uint64 `json:"DealPrice" lightning:"DealPrice"`
}

func (p *Product) PrintInfo(log *golog.Logger) {
	log.Debug("Get Product: ID:" + strconv.FormatUint(uint64(p.ID), 10) + " Name: " + p.ProductName)
}

func (p *Product) FormatProductPrice() string {
	price := uint64(9999999)
	if p.ProductPrice != 0 {
		price = p.ProductPrice
	}
	intString := strconv.FormatUint(price, 10)
	length := len(intString)
	if length < 3 {
		intString = fmt.Sprintf("%02s", intString)
	}
	return intString[:length-2] + "." + intString[length-2:]
}

func (p *Product) FormatDealPrice() string {
	price := uint64(9999999)
	if p.DealPrice != 0 {
		price = p.DealPrice
	}
	intString := strconv.FormatUint(price, 10)
	length := len(intString)
	if length < 3 {
		intString = fmt.Sprintf("%02s", intString)
	}
	return intString[:length-2] + "." + intString[length-2:]
}

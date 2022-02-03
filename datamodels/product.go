package datamodels

import (
	"github.com/kataras/golog"
	"strconv"
)

type Product struct {
	ID           int64  `json:"id" sql:"ID" lightning:"id"`
	ProductName  string `json:"ProductName" sql:"ProductName" lightning:"ProductName"`
	ProductNum   int64  `json:"ProductNum" sql:"ProductNum" lightning:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"ProductImage" lightning:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"ProductUrl" lightning:"ProductUrl"`
}

func (p *Product) PrintInfo(log *golog.Logger) {
	log.Debug("Get Product: ID:" + strconv.FormatInt(p.ID, 10) + " Name: " + p.ProductName)
}

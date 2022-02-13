package main

import (
	"jzmall/common"
	"jzmall/datamodels"
)

func main() {
	data := map[string]string{
		"ID":           "1",
		"productName":  "imooc",
		"productNum":   "2",
		"productImage": "123.jpg",
		"productUrl":   "http://url",
	}
	product := &datamodels.Product{}

	common.DataToStructByTagSql(data, product)
}

package datamodels

type Product struct {
	ID           int64  `json:"id" sql: "ID" lightning:"id"`
	ProductName  string `json:"ProductName" sql:"productName" lightning:"ProductName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" lightning:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" lightning:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" lightning:"ProductUrl"`
}

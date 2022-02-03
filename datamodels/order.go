package datamodels

type Order struct {
	ID          int64 `sql:"ID"`
	UserId      int64 `sql:"UserId"`
	ProductId   int64 `sql:"ProductId"`
	OrderStatus int64 `sql:"OrderStatus"`
}

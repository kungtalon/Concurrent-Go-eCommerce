package datamodels

type Message struct {
	ProductId uint
	UserId    uint
}

func NewMessage(userId uint, productId uint) *Message {
	return &Message{productId, userId}
}

package services

import (
	"products/datamodels"
	"products/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
}

type OrderService struct {
	OrderRepository repositories.IOrderRepository
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService {
	return &OrderService{OrderRepository: repository}
}

func (o *OrderService) GetOrderByID(orderId int64) (order *datamodels.Order, err error) {
	return o.OrderRepository.SelectByKey(orderId)
}

func (o *OrderService) DeleteOrderByID(orderId int64) bool {
	return o.OrderRepository.Delete(orderId)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error {
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() ([]*datamodels.Order, error) {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.OrderRepository.SelectAllWithInfo()
}

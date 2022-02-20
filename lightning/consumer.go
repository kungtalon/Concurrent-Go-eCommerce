package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"jzmall/common"
	"jzmall/datamodels"
	"jzmall/distributed"
	"jzmall/repositories"
	"jzmall/services"
	"log"
)

type OrderConsumer struct {
	ProductService services.IProductService
	OrderService   services.IOrderService
}

func (oc *OrderConsumer) DoConsume(dl *amqp.Delivery) {
	message := &datamodels.Message{}
	err := json.Unmarshal(dl.Body, message)
	if err != nil {
		log.Println(err)
	}
	// insert new order
	_, err = oc.OrderService.InsertOrderByMessage(message)
	if err != nil {
		log.Println(err)
	}

	// reduce product number
	err = oc.ProductService.SubProductByOne(message.ProductId)
	if err != nil {
		log.Println(err)
	}

	// ack false means current message has been consumed
	// ack true means all incoming messages have been consumed
	// must be false!
	dl.Ack(false)
	log.Println("One message consumed successfully: " + string(dl.Body))
}

func main() {
	db, err := common.NewMysqlConnGorm()
	if err != nil {
		fmt.Println(err)
	}
	product := repositories.NewProductManager(db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManagerRepository(db)
	orderService := services.NewOrderService(order)
	consumerHandle := OrderConsumer{productService, orderService}

	rabbitmqConsumeSimple := distributed.NewRabbitMQSimple(common.AMQP_QUEUE_NAME)
	rabbitmqConsumeSimple.ConsumeSimple(&consumerHandle)
}

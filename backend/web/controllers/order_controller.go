package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"products/services"
	"strconv"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService
}

func (o *OrderController) Get() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}
}

func (o *OrderController) GetDelete() {
	idString := o.Ctx.URLParam("id")
	orderId, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	success := o.OrderService.DeleteOrderByID(uint(orderId))
	if success {
		o.Ctx.Application().Logger().Debug("Order Deleted: ID: " + idString)
	} else {
		o.Ctx.Application().Logger().Debug("Order Deletion Failed... ID: " + idString)
	}
	o.Ctx.Redirect("/order/all")
}

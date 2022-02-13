package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"jzmall/datamodels"
	"jzmall/services"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
}

func (p *ProductController) GetDetail() mvc.View {
	id, err := strconv.ParseUint(p.Ctx.URLParam("id"), 10, 32)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
		id = 1
	}
	product, err := p.ProductService.GetProductByID(uint(id))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() mvc.View {
	productIdStr := p.Ctx.URLParam("productID")
	userIdStr := p.Ctx.GetCookie("uid")
	productId, err := strconv.ParseUint(productIdStr, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
		return mvc.View{
			Layout: "shared/productLayout.html",
			Name:   "product/result.html",
			Data: iris.Map{
				"product":     &datamodels.Product{},
				"showMessage": "The product could not be found...",
			},
		}
	}
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	product, err := p.ProductService.GetProductByID(uint(productId))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	showMessage := "The product has run out..."
	if product.ProductNum <= 0 {
		return mvc.View{
			Layout: "shared/productLayout.html",
			Name:   "product/result.html",
			Data: iris.Map{
				"product":     product,
				"showMessage": showMessage,
			},
		}
	}

	product.ProductNum -= 1
	err = p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	order := &datamodels.Order{
		UserId:      uint(userId),
		ProductId:   uint(productId),
		OrderStatus: datamodels.OrderSuccess,
	}
	_, err = p.OrderService.InsertOrder(order)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Application().Logger().Debug(fmt.Sprintf("New Order Created: ID=%d", order.ID))
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     order.ID,
			"showMessage": "Your order has successfully been placed!",
		},
	}
}

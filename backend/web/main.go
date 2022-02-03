package main

import (
	"context"
	"products/backend/web/controllers"
	"products/common"
	"products/repositories"
	"products/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	// create iris object
	app := iris.New()
	// set debugger mode
	log := app.Logger()
	log.SetLevel("debug")
	// template
	templates := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(templates)
	// set up template targets
	app.HandleDir("/assets", "./backend/web/assets")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "Error Occurred..."))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// register controllers
	productRepository := repositories.NewProductManager("concurrent-go-product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/products")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	// app start
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}

package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"jzmall/admin/web/controllers"
	"jzmall/common"
	"jzmall/datamodels"
	"jzmall/repositories"
	"jzmall/services"

	"github.com/kataras/iris/v12/mvc"
)

func main() {
	// create iris object
	app := iris.New()
	// set debugger mode
	log := app.Logger()
	log.SetLevel("debug")
	// template
	templates := iris.HTML("./admin/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(templates)
	// set up template targets
	app.HandleDir("/assets", common.CDN_DOMAIN_URL+"/assets")
	//app.HandleDir("/assets", "./admin/web/assets")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "Error Occurred..."))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	gormdb, err := common.NewMysqlConnGorm()
	if err != nil {
		log.Error(err)
	}
	err = gormdb.AutoMigrate(&datamodels.Product{}, &datamodels.Order{})
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// register controllers
	productRepository := repositories.NewProductManager(gormdb)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	// register productService as a dependency of the controller
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManagerRepository(gormdb)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	// app start
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}

// RegisterNewController reduces the code for registering controllers
func RegisterNewController(ctx context.Context, relativePath string, controller interface{}) {

}

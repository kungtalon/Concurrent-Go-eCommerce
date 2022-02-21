package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"gorm.io/gorm"
	"jzmall/common"
	"jzmall/datamodels"
	"jzmall/distributed"
	"jzmall/front/middleware"
	"jzmall/front/web/controllers"
	"jzmall/repositories"
	"jzmall/services"
)

func main() {
	// create iris object
	app := iris.New()
	// set debugger mode
	log := app.Logger()
	log.SetLevel("debug")
	// template
	templates := iris.HTML("./front/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(templates)
	// set up template targets
	//app.HandleDir("/public", "./front/web/public")
	app.HandleDir("/public", common.CDN_DOMAIN_URL+"/public")
	//app.HandleDir("/html", "./front/web/htmlProductShow")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "An error Occurred..."))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	gormdb, err := common.NewMysqlConnGorm()
	if err != nil {
		log.Error(err)
	}
	err = gormdb.AutoMigrate(&datamodels.User{})
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rabbitmq := distributed.NewRabbitMQSimple(common.AMQP_QUEUE_NAME)

	RegisterControllers(app, gormdb, ctx, rabbitmq)

	// app start
	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}

// register controllers
func RegisterControllers(app *iris.Application, db *gorm.DB, ctx context.Context, rabbitmq *distributed.RabbitMQ) {
	// user controller
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	orderRepository := repositories.NewOrderManagerRepository(db)
	orderService := services.NewOrderService(orderRepository)

	// product controller
	productRepository := repositories.NewProductManager(db)
	productService := services.NewProductService(productRepository)
	productApp := app.Party("/product")
	productApp.Use(middleware.AuthConProduct)
	productPro := mvc.New(productApp)
	productPro.Register(productService, orderService, ctx, rabbitmq)
	productPro.Handle(new(controllers.ProductController))

}

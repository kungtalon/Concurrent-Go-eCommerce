package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"gorm.io/gorm"
	"jzmall/common"
	"jzmall/datamodels"
	"jzmall/frontend/middleware"
	"jzmall/frontend/web/controllers"
	"jzmall/repositories"
	"jzmall/services"
	"time"

	"github.com/kataras/iris/v12/mvc"
)

func main() {
	// create iris object
	app := iris.New()
	// set debugger mode
	log := app.Logger()
	log.SetLevel("debug")
	// template
	templates := iris.HTML("./frontend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(templates)
	// set up template targets
	app.HandleDir("/public", "./frontend/web/public")
	app.HandleDir("/html", "./frontend/web/htmlProductShow")
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

	sess := sessions.New(sessions.Config{
		Cookie:  "helloworld",
		Expires: 60 * time.Minute,
	})

	RegisterControllers(app, gormdb, ctx, sess)

	// app start
	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}

// register controllers
func RegisterControllers(app *iris.Application, db *gorm.DB, ctx context.Context, session *sessions.Sessions) {
	// user controller
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx, session.Start)
	userPro.Handle(new(controllers.UserController))

	orderRepository := repositories.NewOrderManagerRepository(db)
	orderService := services.NewOrderService(orderRepository)

	// product controller
	productRepository := repositories.NewProductManager(db)
	productService := services.NewProductService(productRepository)
	productApp := app.Party("/product")
	productApp.Use(middleware.AuthConProduct)
	productPro := mvc.New(productApp)
	productPro.Register(productService, orderService, ctx, session.Start)
	productPro.Handle(new(controllers.ProductController))

}

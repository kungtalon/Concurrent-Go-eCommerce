package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	// create iris object
	app := iris.New()

	// template
	//templates := iris.HTML("./front/web/views", ".html").Layout("shared/layout.html").Reload(true)
	//app.RegisterView(templates)
	app.HandleDir("/public", "./front/web/public")
	app.HandleDir("/html", "./front/web/htmlProductShow")
	// app start
	app.Run(
		iris.Addr("0.0.0.0:80"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}

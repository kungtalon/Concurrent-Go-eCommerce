package middleware

import "github.com/kataras/iris/v12"

func AuthConProduct(ctx iris.Context) {
	uid := ctx.GetCookie("uid")
	if uid == "" {
		ctx.Application().Logger().Debug("User has not logged in!")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("User has logged in!")
	ctx.Next()
}

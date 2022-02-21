package middleware

import (
	"github.com/kataras/iris/v12"
	"jzmall/common"
)

func AuthConProduct(ctx iris.Context) {
	uid := ctx.GetCookie("uid")
	uidStr := ctx.GetCookie("sign")
	if uid == "" || uidStr == "" {
		ctx.Application().Logger().Debug("User has not logged in!")
		ctx.Redirect("/user/login")
		return
	}
	decoded, err := common.DePwdCode(uidStr)
	if err != nil {
		ctx.Application().Logger().Debug("Error when decoding encoded userid...")
		ctx.Application().Logger().Debug(err)
		ctx.Redirect("/user/login")
		return
	}
	if string(decoded) != uid {
		ctx.Application().Logger().Debug("Invalid user information found in cookies! Logged out...")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("User has logged in!")
	ctx.Next()
}

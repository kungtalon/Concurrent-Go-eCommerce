package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"jzmall/common"
	"jzmall/datamodels"
	"jzmall/services"
	"strconv"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
	Session *sessions.Session
}

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

func (c *UserController) PostRegister() {
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)

	user := &datamodels.User{
		UserName:     userName,
		NickName:     nickName,
		HashPassword: password,
	}

	_, err := c.Service.AddUser(user)
	if err != nil {
		c.Ctx.Application().Logger().Debug(err)
		c.Ctx.Redirect("/user/error")
		return
	}

	c.Ctx.Redirect("/user/login")
	return
}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (c *UserController) PostLogin() mvc.Response {
	var (
		userName = c.Ctx.FormValue("userName")
		passWord = c.Ctx.FormValue("password")
	)

	// validate userid and password
	user, success := c.Service.IsPwdSuccess(userName, passWord)
	if !success {
		c.Ctx.Application().Logger().Debug("Wrong Password for User " + userName)
		return mvc.Response{
			Path: "/user/login",
		}
	}

	// write user id to cookie
	common.GlobalCookie(c.Ctx, "uid", strconv.FormatUint(uint64(user.ID), 10))
	uidByte := []byte(strconv.FormatUint(uint64(user.ID), 10))
	uidStr, err := common.EnPwdCode(uidByte)
	if err != nil {
		c.Ctx.Application().Logger().Debug(err)
	}
	common.GlobalCookie(c.Ctx, "sign", uidStr)
	return mvc.Response{
		Path: "/product",
	}
}

func (c *UserController) GetPopulate() {
	for i := 5; i <= 200; i++ {
		nickName := "guest" + strconv.Itoa(i)
		userName := "guest" + strconv.Itoa(i)
		passWord := "guest" + strconv.Itoa(i)
		newUser := &datamodels.User{NickName: nickName, UserName: userName, HashPassword: passWord}
		_, err := c.Service.AddUser(newUser)
		if err != nil {
			c.Ctx.Application().Logger().Debug(err)
		}
	}

	c.Ctx.Redirect("/user/login")
}

package datamodels

import "gorm.io/gorm"

type User struct {
	gorm.Model
	NickName     string `json:"NickName" form:"nickName" sql:"NickName"`
	UserName     string `json:"UserName" form:"userName" sql:"NickName"`
	HashPassword string `json:"-" form:"passWord" sql:"passWord"`
}

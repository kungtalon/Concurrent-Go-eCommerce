package services

import (
	"golang.org/x/crypto/bcrypt"
	"jzmall/datamodels"
	"jzmall/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, success bool)
	AddUser(user *datamodels.User) (uint, error)
}

type UserService struct {
	UserRepository repositories.IUser
}

func NewUserService(repository repositories.IUser) IUserService {
	return &UserService{UserRepository: repository}
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, success bool) {
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		return
	}

	success, _ = ValidatePassword(pwd, user.HashPassword)
	if !success {
		return &datamodels.User{}, false
	}
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId uint, err error) {
	pwdBytes, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return userId, err
	}
	user.HashPassword = string(pwdBytes)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPwd), bcrypt.DefaultCost)
}

func ValidatePassword(userPwd string, hashed string) (success bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPwd)); err != nil {
		return false, err
	}
	return true, nil
}

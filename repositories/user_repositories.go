package repositories

import (
	"errors"
	"gorm.io/gorm"
	"jzmall/common"
	"jzmall/datamodels"
)

type IUser interface {
	Conn() error
	Select(userName string) (*datamodels.User, error)
	Insert(user *datamodels.User) (uint, error)
}

type UserManagerRepository struct {
	mysqlConn *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUser {
	return &UserManagerRepository{mysqlConn: db}
}

var _ IUser = (*UserManagerRepository)(nil)

func (u *UserManagerRepository) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, err := common.NewMysqlConnGorm()
		if err != nil {
			return err
		}
		u.mysqlConn = mysql
		err = u.mysqlConn.AutoMigrate(&datamodels.User{})
		if err != nil {
			return err
		}
	}
	return
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("Input must not be empty!")
	}
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	user = &datamodels.User{}
	err = u.mysqlConn.Where("user_name=?", userName).First(&user).Error
	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId uint, err error) {
	if err = u.Conn(); err != nil {
		return
	}

	record := &datamodels.User{}
	result := u.mysqlConn.Model(&record).Where("user_name=?", user.UserName).Limit(1).Find(&record)
	if result.Error != nil {
		return
	}
	if result.RowsAffected > 0 {
		return 0, errors.New("User name already exists!")
	}
	err = u.mysqlConn.Create(&user).Error
	return user.ID, err
}

func (u *UserManagerRepository) SelectById(userId uint) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	user = &datamodels.User{}
	err = u.mysqlConn.Where("ID=?", userId).First(&user).Error

	return
}

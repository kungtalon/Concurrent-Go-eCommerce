package repositories

import (
	"gorm.io/gorm"
	"jzmall/common"
	"jzmall/datamodels"
)

// develop the interface
// implement the interface

type IProduct interface {
	Conn() error
	Insert(*datamodels.Product) (uint, error)
	Delete(uint) bool
	Update(*datamodels.Product) error
	SelectByKey(uint) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductManager struct {
	mysqlConn *gorm.DB
}

func NewProductManager(db *gorm.DB) IProduct {
	return &ProductManager{mysqlConn: db}
}

func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConnGorm()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
		err = p.mysqlConn.AutoMigrate(&datamodels.Product{})
		if err != nil {
			return err
		}
	}
	return
}

func (p *ProductManager) Insert(product *datamodels.Product) (productID uint, err error) {
	// check connection
	if err = p.Conn(); err != nil {
		return
	}

	result := p.mysqlConn.Create(&product)

	return product.ID, result.Error
}

func (p *ProductManager) Delete(productID uint) bool {
	// check connection
	if err := p.Conn(); err != nil {
		return false
	}

	result := p.mysqlConn.Delete(&datamodels.Product{}, productID)
	if result.Error != nil {
		return false
	}
	return true
}

func (p *ProductManager) Update(product *datamodels.Product) (err error) {
	if err = p.Conn(); err != nil {
		return err
	}

	result := p.mysqlConn.Model(&datamodels.Product{}).Where("ID=?", product.ID).Updates(product)
	return result.Error
}

// select a product by its ID
func (p *ProductManager) SelectByKey(productID uint) (productResult *datamodels.Product, err error) {
	// check connection
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}

	productResult = &datamodels.Product{}
	result := p.mysqlConn.First(&productResult, productID)
	err = result.Error

	return
}

func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, err error) {
	if err = p.Conn(); err != nil {
		return nil, err
	}

	results := p.mysqlConn.Find(&productArray)
	err = results.Error

	return
}

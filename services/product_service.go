package services

import (
	"jzmall/datamodels"
	"jzmall/repositories"
)

type IProductService interface {
	GetProductByID(uint) (*datamodels.Product, error)
	GetAllProducts() ([]*datamodels.Product, error)
	DeleteProductByID(uint) bool
	InsertProduct(*datamodels.Product) (uint, error)
	UpdateProduct(*datamodels.Product) error
}

type ProductService struct {
	productRepository repositories.IProduct
}

func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{productRepository: repository}
}

func (p *ProductService) GetProductByID(productID uint) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(productID)
}

func (p *ProductService) GetAllProducts() ([]*datamodels.Product, error) {
	return p.productRepository.SelectAll()
}

func (p *ProductService) DeleteProductByID(productID uint) bool {
	return p.productRepository.Delete(productID)
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (uint, error) {
	return p.productRepository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	return p.productRepository.Update(product)
}

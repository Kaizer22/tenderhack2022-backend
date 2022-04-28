package repository

import (
	"context"
	"main/model/entity"
)

const (
	TagInsPr     = "INSERT PRODUCTS"
	TagUpdPr     = "UPDATE PRODUCTS"
	TagDelPr     = "DELETE PRODUCTS"
	TagGetAllPr  = "GET ALL PRODUCTS"
	TagGetPrById = "GET PRODUCTS BY ID"
	TagGetPrByQ  = "GET PRODUCTS BY QUERY"
)

type ProductRepository interface {
	InsertProducts(ctx context.Context, products []*entity.Product) error
	UpdateProduct(ctx context.Context, id int, data entity.ProductData) error
	DeleteProducts(ctx context.Context, id int) error
	GetAllProducts(ctx context.Context) ([]*entity.Product, error)
	GetProductById(ctx context.Context, id int) (*entity.Product, error)
	GetProductsByQuery(ctx context.Context, q string) ([]*entity.Product, error)
}

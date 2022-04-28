package repository

import (
	"context"
	"main/model/entity"
)

const (
	TagInsCat     = "INSERT CATEGORY"
	TagUpdCat     = "UPDATE CATEGORY"
	TagDelCat     = "DELETE CATEGORY"
	TagGetAllCats = "GET ALL CATEGORIES"
	TagGetCatById = "GET CATEGORY BY ID"
)

type CategoryRepository interface {
	InsertCategories(ctx context.Context, products []*entity.Category) error
	UpdateCategory(ctx context.Context, id int, data entity.CategoryData) error
	DeleteCategory(ctx context.Context, id int) error
	GetAllCategories(ctx context.Context) ([]*entity.Category, error)
	GetCategoryById(ctx context.Context, id int) (*entity.Category, error)
}

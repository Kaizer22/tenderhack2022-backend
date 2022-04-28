package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/model/entity"
	"main/repository"
	"main/utils"
	"net/http"
	"strconv"
)

type CategoryController struct {
	CategoryRepo repository.CategoryRepository
	ctx          context.Context
}

// NewCategoryController example
func NewCategoryController(ctx context.Context, repo repository.CategoryRepository) *CategoryController {
	return &CategoryController{
		CategoryRepo: repo,
		ctx:          ctx,
	}
}

// AddCategory godoc
// @Summary            Add new category
// @Description    Add new category and get entity with ID in a response
// @Tags                          categories
// @Accept                        json
// @Produce                       json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               category   body            entity.CategoryData  true  "Category info"
// @Success             201             {string}  string                        "New category successfully added"
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /api/v1/categories [post]
func (c *CategoryController) AddCategory(ctx *gin.Context) {
	var p entity.CategoryData
	if err := ctx.ShouldBindJSON(&p); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	//TODO add struct validation
	category := entity.Category{
		Name:        p.Name,
		Description: p.Description,
	}

	//TODO deal with contexts correctly
	c2 := context.Background()
	err := c.CategoryRepo.InsertCategories(c2, []*entity.Category{
		&category,
	})
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	//TODO return struct with ID
	ctx.JSON(http.StatusOK, "Category created")
}

// ListCategories godoc
// @Summary            Get categories
// @Description    Returns all the categories in system
// @Tags                       categories
// @Accept                     json
// @Produce                    json
// @Param        Authorization  header    string  true  "Authentication header"
// @Success             200              {array}         entity.Category
// @Failure        400         {object}           utils.HTTPError
// @Failure        404         {object}           utils.HTTPError
// @Failure        500         {object}           utils.HTTPError
// @Router                     /api/v1/categories [get]
func (c *CategoryController) ListCategories(ctx *gin.Context) {

	c2 := context.Background()
	if categories, err := c.CategoryRepo.GetAllCategories(c2); err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		ctx.JSON(http.StatusOK, categories)
		return
	}
}

// GetCategory godoc
// @Summary            Get category by ID
// @Description    Returns category by ID
// @Tags                     categories
// @Accept                   json
// @Produce                  json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               id             path        int  true  "Category ID"
// @Success             200            {object}  entity.Category
// @Failure        400       {object}  utils.HTTPError
// @Failure        404       {object}  utils.HTTPError
// @Failure        500       {object}  utils.HTTPError
// @Router                   /api/v1/categories/{id} [get]
func (c *CategoryController) GetCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	iID, err := strconv.Atoi(id)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	//TODO deal with contexts correctly
	c2 := context.Background()
	product, err := c.CategoryRepo.GetCategoryById(c2, iID)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// PutCategory godoc
// @Summary            Edit category
// @Description    Edit existing category
// @Tags                      categories
// @Accept                    json
// @Produce                   json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               id                        path                        int         true  "Category ID"
// @Param               data            body              entity.CategoryData  true  "Category entity"
// @Success             200             {string}  string  "Category updated"
// @Failure        400        {object}            utils.HTTPError
// @Failure        404        {object}            utils.HTTPError
// @Failure        500        {object}            utils.HTTPError
// @Router                    /api/v1/categories/{id} [put]
func (c *CategoryController) PutCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	iID, err := strconv.Atoi(id)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	var p entity.CategoryData
	if err := ctx.ShouldBindJSON(&p); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	//TODO add struct validation
	product := entity.CategoryData{
		Name:        p.Name,
		Description: p.Description,
	}

	//TODO deal with contexts correctly
	c2 := context.Background()
	err = c.CategoryRepo.UpdateCategory(c2, iID, product)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "Category updated")
}

// DeleteCategory godoc
// @Summary            Delete category
// @Description    Delete selected category
// @Tags                          categories
// @Accept                        json
// @Produce                       json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               id                  path    int    true  "Category ID"
// @Success             200       {string}  string  "Category deleted"
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /api/v1/categories/{id} [delete]
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	productId := ctx.Param("id")
	c2 := context.Background()
	iID, err := strconv.Atoi(productId)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	err = c.CategoryRepo.DeleteCategory(c2, iID)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "Category deleted")
}

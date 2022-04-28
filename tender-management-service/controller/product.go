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

// Example controller for product entity
type ProductController struct {
	ProductRepo repository.ProductRepository
	ctx         context.Context
}

// NewProductController example
func NewProductController(ctx context.Context, repo repository.ProductRepository) *ProductController {
	return &ProductController{
		ProductRepo: repo,
		ctx:         ctx,
	}
}

// AddProduct godoc
// @Summary            Add new product
// @Description    Add new product and get entity with ID in a response
// @Tags                          products
// @Accept                        json
// @Produce                       json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               product   body            entity.ProductData  true  "Product info"
// @Success             201             {string}  string                        "New product successfully added"
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /api/v1/products [post]
func (c *ProductController) AddProduct(ctx *gin.Context) {
	var p entity.ProductData
	if err := ctx.ShouldBindJSON(&p); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	//TODO add struct validation
	product := entity.Product{
		Name:        p.Name,
		Description: p.Description,
		CategoryId:  p.CategoryId,
	}

	//TODO deal with contexts correctly
	c2 := context.Background()
	err := c.ProductRepo.InsertProducts(c2, []*entity.Product{
		&product,
	})
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	//TODO return struct with ID
	ctx.JSON(http.StatusOK, "Product created")
}

// ListProducts godoc
// @Summary            Get products
// @Description    Returns all the products in system or products filtered using query
// @Tags                       products
// @Accept                     json
// @Produce                    json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               query  query     string   false  "search substring in name, description or category"
// @Success             201              {array}         entity.Product
// @Failure        400         {object}           utils.HTTPError
// @Failure        404         {object}           utils.HTTPError
// @Failure        500         {object}           utils.HTTPError
// @Router                     /api/v1/products [get]
func (c *ProductController) ListProducts(ctx *gin.Context) {
	query := ctx.Request.URL.Query().Get("query")
	c2 := context.Background()
	if len(query) > 0 {
		if products, err := c.ProductRepo.GetProductsByQuery(c2, query); err != nil {
			utils.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			ctx.JSON(http.StatusOK, products)
			return
		}
	} else {
		if products, err := c.ProductRepo.GetAllProducts(c2); err != nil {
			utils.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			ctx.JSON(http.StatusOK, products)
			return
		}
	}
}

// GetProduct godoc
// @Summary            Get product by ID
// @Description    Returns product by ID
// @Tags                     products
// @Accept                   json
// @Produce                  json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               id             path        int  true  "Product ID"
// @Success             200            {object}  entity.Product
// @Failure        400       {object}  utils.HTTPError
// @Failure        404       {object}  utils.HTTPError
// @Failure        500       {object}  utils.HTTPError
// @Router                   /api/v1/products/{id} [get]
func (c *ProductController) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	iID, err := strconv.Atoi(id)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	//TODO deal with contexts correctly
	c2 := context.Background()
	product, err := c.ProductRepo.GetProductById(c2, iID)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// PutProduct godoc
// @Summary            Edit product
// @Description    Edit existing product
// @Tags                      products
// @Accept                    json
// @Produce                   json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               id                        path                        int         true  "Product ID"
// @Param               data            body              entity.ProductData  true  "Product entity"
// @Success             200             {string}  string  "Product updated"
// @Failure        400        {object}            utils.HTTPError
// @Failure        404        {object}            utils.HTTPError
// @Failure        500        {object}            utils.HTTPError
// @Router                    /api/v1/products/{id} [put]
func (c *ProductController) PutProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	iID, err := strconv.Atoi(id)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	var p entity.ProductData
	if err := ctx.ShouldBindJSON(&p); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	//TODO add struct validation
	product := entity.ProductData{
		Name:        p.Name,
		Description: p.Description,
		CategoryId:  p.CategoryId,
	}

	//TODO deal with contexts correctly
	c2 := context.Background()
	err = c.ProductRepo.UpdateProduct(c2, iID, product)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "Product updated")
}

// DeleteProduct godoc
// @Summary            Delete product
// @Description    Delete selected product
// @Tags                          products
// @Accept                        json
// @Produce                       json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               id                  path    int    true  "Account ID"
// @Success             200       {string}  string  "Product deleted"
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /api/v1/products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	productId := ctx.Param("id")
	c2 := context.Background()
	iID, err := strconv.Atoi(productId)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	err = c.ProductRepo.DeleteProducts(c2, iID)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "Product deleted")
}

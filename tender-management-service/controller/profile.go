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
type ProfileController struct {
	ProfileRepo repository.ProfileRepository
	ctx         context.Context
}

// NewProfileController example
func NewProfileController(ctx context.Context, repo repository.ProfileRepository) *ProfileController {
	return &ProfileController{
		ProfileRepo: repo,
		ctx:         ctx,
	}
}

// GetProfile godoc
// @Summary            Get profile by ID
// @Description    Returns profile by ID
// @Tags                     profiles
// @Accept                   json
// @Produce                  json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               id             path        int  true  "Profile ID"
// @Success             200            {object}  entity.Profile
// @Failure        400       {object}  utils.HTTPError
// @Failure        404       {object}  utils.HTTPError
// @Failure        500       {object}  utils.HTTPError
// @Router                   /api/v1/profiles/{id} [get]
func (c *ProfileController) GetProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	iID, err := strconv.Atoi(id)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	//TODO deal with contexts correctly
	c2 := context.Background()
	profile, err := c.ProfileRepo.GetProfileById(c2, int64(iID))
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary            Edit profile
// @Description    Edit existing profile
// @Tags                      profiles
// @Accept                    json
// @Produce                   json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               id                        path                        int         true  "Profile ID"
// @Param               data            body              entity.ProfileData  true  "Profile entity"
// @Success             200             {string}  string  "Profile updated"
// @Failure        400        {object}            utils.HTTPError
// @Failure        404        {object}            utils.HTTPError
// @Failure        500        {object}            utils.HTTPError
// @Router                    /api/v1/profiles/{id} [put]
func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	iID, err := strconv.Atoi(id)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	var p entity.ProfileData
	if err := ctx.ShouldBindJSON(&p); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	//TODO deal with contexts correctly
	c2 := context.Background()
	err = c.ProfileRepo.UpdateProfile(c2, int64(iID), &p)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "Profile updated")
}

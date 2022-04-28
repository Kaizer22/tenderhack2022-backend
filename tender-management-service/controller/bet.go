package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/model/entity"
	"main/model/response"
	"main/repository"
	"main/utils"
	"net/http"
	"strconv"
)

type BetController struct {
	BetRepo repository.BetRepository
	ctx     context.Context
}

// NewProductController example
func NewBetController(ctx context.Context, repo repository.BetRepository) *BetController {
	return &BetController{
		BetRepo: repo,
		ctx:     ctx,
	}
}

// MakeBet godoc
// @Summary            Make new bet
// @Description    Make new bet in a quotation session
// @Tags                          bets
// @Accept                        json
// @Produce                       json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               bet   body            entity.BetData  true  "Bet info"
// @Success             201             {string}  string                        "New bet successfully made"
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /api/v1/bets [post]
func (c *BetController) MakeBet(ctx *gin.Context) {
	var p entity.BetData
	if err := ctx.ShouldBindJSON(&p); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	c2 := context.Background()
	id, err := c.BetRepo.MakeBet(c2, p)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, response.BetMade{ID: id})
}

// GetBetsBySessionId godoc
// @Summary            Get bets by session Id
// @Description    Get bets history for a selected session
// @Tags                          bets
// @Accept                        json
// @Produce                       json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               sessionId   path            int true  "Session id"
// @Success             200             {array}  entity.Bet
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /api/v1/bets/{sessionId} [get]
func (c *BetController) GetBetsBySessionId(ctx *gin.Context) {
	id := ctx.Param("sessionId")
	iID, err := strconv.Atoi(id)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	//TODO deal with contexts correctly
	c2 := context.Background()
	bets, err := c.BetRepo.GetBetBySessionId(c2, int64(iID))
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, bets)
}

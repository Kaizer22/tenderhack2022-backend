package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/logging"
	"main/model/entity"
	"main/repository"
	"main/service"
	"main/utils"
	"net/http"
	"strconv"
)

// Controller for strategy
type StrategyController struct {
	strategySrv service.StrategyService
	ctx         context.Context
}

// NewStrategyController example
func NewStrategyController(ctx context.Context, repo repository.QuotationSessionRepository,
	betRepo repository.BetRepository) *StrategyController {

	srv := service.NewStrategyService(ctx, repo, betRepo)
	return &StrategyController{
		ctx:         ctx,
		strategySrv: srv,
	}
}

// RunStrategy godoc
// @Summary            Run strategy
// @Description    Run selected strategy
// @Tags                      strategies
// @Accept                    json
// @Produce                   json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               data            body      entity.StrategyParams true  "Strategy params"
// @Success             200             {string}  string  "Strategy launched"
// @Failure        400        {object}            utils.HTTPError
// @Failure        404        {object}            utils.HTTPError
// @Failure        500        {object}            utils.HTTPError
// @Router                    /api/v1/strategies/run [post]
func (c StrategyController) RunStrategy(ctx *gin.Context) {
	var p entity.StrategyParams
	if err := ctx.ShouldBindJSON(&p); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	err := c.strategySrv.RunStrategyRunner(p)
	if err != nil {
		logging.ErrorFormat("Error launching strategy runner: %s", err)
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "Runner successfully launched.")
}

// StopStrategy godoc
// @Summary            Stop user strategy
// @Description    Stops selected strategy
// @Tags                      strategies
// @Accept                    json
// @Produce                   json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               sessionId            path      int true  "Session ID param"
// @Param               userId            path      int true  "User ID param"
// @Success             200             {string}  string  "Strategy stopped"
// @Failure        400        {object}            utils.HTTPError
// @Failure        404        {object}            utils.HTTPError
// @Failure        500        {object}            utils.HTTPError
// @Router                    /api/v1/strategies/{sessionId}/{userId} [post]
func (c StrategyController) StopStrategy(ctx *gin.Context) {
	sessionId := ctx.Param("sessionId")
	isessionID, err := strconv.Atoi(sessionId)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	userId := ctx.Param("userId")
	iuserID, err := strconv.Atoi(userId)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	err = c.strategySrv.StopRunner(int64(isessionID), int64(iuserID))
	if err != nil {
		logging.ErrorFormat("Error stopping strategy runner: %s", err)
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "Runner successfully stopped.")
}

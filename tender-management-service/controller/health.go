package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/model/response"
	"main/utils"
	"net/http"
)

type HealthCheckComponent interface {
	IsConnected() (bool, error)
	Description() string
}

type HealthCheckController struct {
	components []HealthCheckComponent
	ctx        context.Context
}

func NewHealthCheckController(ctx context.Context, components ...HealthCheckComponent) HealthCheckController {
	return HealthCheckController{
		components: components,
		ctx:        ctx,
	}
}

// GetHealthStatus godoc
// @Summary            Get health
// @Description    Health status for services' components
// @Tags                          health
// @Accept                        json
// @Produce                       json
// @Success             200             {array}  response.HealthStatus
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /health [get]
func (c HealthCheckController) GetHealthStatus(ctx *gin.Context) {
	var stat []response.HealthStatus

	stat = append(stat, response.HealthStatus{
		Status: response.StatusUp,
		Info:   utils.GetEnv(utils.ServiceNameEnvKey, "golang-template-microservice"),
	})

	status := response.StatusUnknown
	for _, com := range c.components {
		s, err := com.IsConnected()
		if err != nil {

		}
		if s {
			status = response.StatusUp
		} else {
			status = response.StatusDown
		}

		stat = append(stat, response.HealthStatus{
			Status: status,
			Info:   com.Description(),
		})
	}

	ctx.JSON(http.StatusOK, stat)
}

// GetServiceVersion godoc
// @Summary            Get version
// @Description    Current version of the service
// @Tags                          health
// @Accept                        json
// @Produce                       json
// @Success        200       {string} string
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /health/version [get]
func (c HealthCheckController) GetServiceVersion(ctx *gin.Context) {
	version := utils.GetEnv(utils.AppVersionEnvKey, "v0_dummy")
	ctx.JSON(http.StatusOK, version)
}

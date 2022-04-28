package main

import (
	"context"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/gin-gonic/gin"
	"main/config"
	"main/controller"
	conn "main/db/impl"
	"main/docs"
	"main/logging"
	"main/middleware"
	repo "main/repository/impl"
	"main/utils"

	"github.com/gin-contrib/pprof"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "main/docs"
)

var (
	addr string
)

// @title           Strategy Management Service API
// @version         1.0
// @description     Service for async running of user strategies.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8083
// @BasePath  /api/v1

// @securityDefinitions.basic  JWT

const (
	apiV1 = "/api/v1"
)

// TODO users and basic auth
func main() {
	ctx := context.Background()
	r := gin.Default()

	connection, err := conn.NewPgOrmConnectionProvider()
	if err != nil {
		panic(err)
	}

	sessionRepo := repo.NewPgOrmQuotationSessionRepository(connection.Connection().(*pg.DB))
	betRepo := repo.NewPgOrmBetRepository(ctx, connection.Connection().(*pg.DB))

	strategyC := controller.NewStrategyController(ctx, sessionRepo, betRepo)

	hC := controller.NewHealthCheckController(ctx,
		connection,
	)
	fileAdapter := fileadapter.NewAdapter("config/basic_policy.csv")

	authorized := r.Group("/")
	authorized.Use(gin.Recovery())
	authorized.Use(middleware.TokenAuthMiddleware())
	{
		v1 := r.Group(apiV1)
		{
			strategies := v1.Group("/strategies")
			{
				strategies.POST("/run", middleware.Authorize(config.Strategy, config.Run, fileAdapter),
					strategyC.RunStrategy)
				strategies.POST(":sessionId/:userId", middleware.Authorize(config.Strategy, config.Stop, fileAdapter),
					strategyC.StopStrategy)
			}

		}
		//authorized.POST("/logout", authC.Logout)
		//authorized.POST("/refresh", authC.Refresh)
	}
	health := r.Group("/health")
	{
		health.GET("", hC.GetHealthStatus)
		health.GET("version", hC.GetServiceVersion)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	pprof.Register(r, "/debug/pprof")

	initService()
	logging.InfoFormat("Starting server at %s", addr)
	err = r.Run(addr)
	if err != nil {
		logging.FatalFormat("unable to start server")
		panic(err)
		return
	}

}

func initService() {
	docs.SwaggerInfo.Host = utils.GetEnv(utils.StrategyManagementBaseUrlEnvKey, "localhost:"+
		utils.GetEnv(utils.ListenAddressEnvKey, "8083"))
	docs.SwaggerInfo.BasePath = "/"
	addr = ":" + utils.GetEnv(utils.ListenAddressEnvKey, "8083")
}

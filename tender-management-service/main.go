package main

import (
	"context"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"main/config"
	"main/controller"
	"main/db"
	conn "main/db/impl"
	"main/docs"
	"main/logging"
	"main/middleware"
	repo "main/repository/impl"
	"main/service"
	"main/utils"

	"github.com/gin-contrib/pprof"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "main/docs"
)

var (
	addr string
)

// @title           Tender Management Service API
// @version         1.0
// @description     Service to manage sessions and bets.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.basic  JWT

const (
	apiV1 = "/api/v1"
)

// TODO users and basic auth
func main() {
	ctx := context.Background()
	r := gin.Default()

	//connection, err := conn.NewPgConnectionProvider()
	connection, err := conn.NewPgOrmConnectionProvider()
	if err != nil {
		panic(err)
	}
	err = connection.Migrate(db.PgMigrationsPath)
	if err != nil {
		panic(err)
	}

	// Two possible options: pg ORM with schema init or SQL queries based repository
	// with custom migrations from files
	//productRepo := repo.NewPgProductRepository(connection.Connection())
	//------------------------------------------------------------------------------
	productRepo := repo.NewPgOrmProductRepository(ctx, connection.Connection().(*pg.DB))
	categoryRepo := repo.NewPgOrmCategoryRepository(ctx, connection.Connection().(*pg.DB))
	accountRepo := repo.NewPgOrmAccountRepository(ctx, connection.Connection().(*pg.DB))
	profileRepo := repo.NewPgOrmProfileRepository(ctx, connection.Connection().(*pg.DB))
	sessionRepo := repo.NewPgOrmQuotationSessionRepository(connection.Connection().(*pg.DB))
	betRepo := repo.NewPgOrmBetRepository(ctx, connection.Connection().(*pg.DB))
	pJRepo := repo.NewPgOrmProductJournalRepository(ctx, connection.Connection().(*pg.DB))

	c := controller.NewProductController(ctx, productRepo)
	catC := controller.NewCategoryController(ctx, categoryRepo)
	authC := controller.NewAuthController(ctx, accountRepo)
	prfC := controller.NewProfileController(ctx, profileRepo)
	qsC := controller.NewSessionController(ctx, sessionRepo, pJRepo)
	betC := controller.NewBetController(ctx, betRepo)

	qsSrv := service.NewQuotationSessionService(ctx, sessionRepo)
	err = qsSrv.RunQuotationService()
	if err != nil {
		logging.ErrorFormat("Cannot run quotation session service: %s", err)
	}

	hC := controller.NewHealthCheckController(ctx,
		connection,
	)
	fileAdapter := fileadapter.NewAdapter("config/basic_policy.csv")

	r.POST("/login", authC.Login)
	r.POST("/register", authC.Register)

	authorized := r.Group("/")
	authorized.Use(gin.Recovery())
	authorized.Use(middleware.TokenAuthMiddleware())
	{
		v1 := r.Group(apiV1)
		{
			products := v1.Group("/products")
			{
				products.POST("", middleware.Authorize(config.Product, config.Insert, fileAdapter),
					c.AddProduct)
				products.GET("", middleware.Authorize(config.Product, config.Read, fileAdapter),
					c.ListProducts)
				products.GET(":id", middleware.Authorize(config.Product, config.Read, fileAdapter),
					c.GetProduct)
				products.PUT(":id", middleware.Authorize(config.Product, config.Update, fileAdapter),
					c.PutProduct)
				products.DELETE(":id", middleware.Authorize(config.Product, config.Delete, fileAdapter),
					c.DeleteProduct)
			}
			categories := v1.Group("/categories")
			{
				categories.POST("", middleware.Authorize(config.Category, config.Insert, fileAdapter),
					catC.AddCategory)
				categories.GET("", middleware.Authorize(config.Category, config.Read, fileAdapter),
					catC.ListCategories)
				categories.GET(":id", middleware.Authorize(config.Category, config.Read, fileAdapter),
					catC.GetCategory)
				categories.PUT(":id", middleware.Authorize(config.Category, config.Update, fileAdapter),
					catC.PutCategory)
				categories.DELETE(":id", middleware.Authorize(config.Category, config.Delete, fileAdapter),
					catC.DeleteCategory)
			}
			sessions := v1.Group("/sessions")
			{
				sessions.GET("", middleware.Authorize(config.Session, config.Read, fileAdapter),
					qsC.GetAllSessions)
				sessions.GET(":id", middleware.Authorize(config.Session, config.Read, fileAdapter),
					qsC.GetSessionById)
				sessions.POST("", middleware.Authorize(config.Session, config.Insert, fileAdapter),
					qsC.NewQuotationSession)
				sessions.PUT("")
				sessions.DELETE("")
			}
			profiles := v1.Group("/profiles")
			{
				profiles.GET(":id", middleware.Authorize(config.Profile, config.Read, fileAdapter),
					prfC.GetProfile)
				profiles.PUT(":id", middleware.Authorize(config.Profile, config.Update, fileAdapter),
					prfC.UpdateProfile)
			}

			bets := v1.Group("/bets")
			{
				bets.POST("", middleware.Authorize(config.Bet, config.Insert, fileAdapter),
					betC.MakeBet)
				bets.GET(":sessionId", middleware.Authorize(config.Bet, config.Read, fileAdapter),
					betC.GetBetsBySessionId)
			}

		}
		authorized.POST("/logout", authC.Logout)
		authorized.POST("/refresh", authC.Refresh)
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
	docs.SwaggerInfo.Host = utils.GetEnv(utils.TenderManagementBaseUrlEnvKey, "localhost:"+
		utils.GetEnv(utils.ListenAddressEnvKey, "8080"))
	docs.SwaggerInfo.BasePath = "/"
	addr = ":" + utils.GetEnv(utils.ListenAddressEnvKey, "8080")
}

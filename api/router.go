package api

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	sc "github.com/Manish-Mehta/tigerhall/internal/controller/sight"
	tc "github.com/Manish-Mehta/tigerhall/internal/controller/tiger"
	uc "github.com/Manish-Mehta/tigerhall/internal/controller/user"
	"github.com/Manish-Mehta/tigerhall/internal/middleware"
	ss "github.com/Manish-Mehta/tigerhall/internal/service/sight"
	ts "github.com/Manish-Mehta/tigerhall/internal/service/tiger"
	us "github.com/Manish-Mehta/tigerhall/internal/service/user"
	"github.com/Manish-Mehta/tigerhall/model/datastore"
	"github.com/Manish-Mehta/tigerhall/pkg/db"
)

// @title 	Tiger service API
// @version	1.0
// @description Tiger service API for tiger management and recording system

// @host 	localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Description for what is this security definition being used
func SetupRouter(engine *gin.Engine) {
	log.Println("Initializing Routes")
	dBClient := db.GetDBClient().GetClient()
	engine.MaxMultipartMemory = 8 << 20 // 8 MiB

	apiRouter := engine.Group("/api/v1")
	apiRouter.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	{
		apiRouter.HEAD("/health", HealthApi)
		apiRouter.GET("/health", HealthApi)

		// User Router
		{
			userRouter := apiRouter.Group("/user")

			userStore := datastore.NewUserStore(dBClient)
			userService := us.NewUserService(userStore)
			userController := uc.NewUserController(userService)
			{
				userRouter.POST("", userController.Signup)
				userRouter.POST("/login", userController.Login)
				userRouter.GET("/refresh", middleware.AuthMiddleware, userController.Refresh)

			}
		}
		// Tiger Router
		tigerStore := datastore.NewTigerStore(dBClient)
		{
			tigerRouter := apiRouter.Group("/tiger")
			tigerService := ts.NewTigerService(tigerStore)

			tigerController := tc.NewTigerController(tigerService)
			{
				tigerRouter.POST("", middleware.AuthMiddleware, tigerController.Create)
				tigerRouter.GET("", tigerController.List)
			}
		}
		// Sight Router
		{
			sightRouter := apiRouter.Group("/sight")
			sightStore := datastore.NewSightStore(dBClient)
			sightService := ss.NewSightService(sightStore, tigerStore)

			sightController := sc.NewSightController(sightService)
			{
				sightRouter.POST("", middleware.AuthMiddleware, sightController.Create)
				sightRouter.GET("", sightController.List)
			}
		}
	}
}

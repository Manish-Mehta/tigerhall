package controller

import (
	"log"

	sc "github.com/Manish-Mehta/tigerhall/internal/controller/sight"
	tc "github.com/Manish-Mehta/tigerhall/internal/controller/tiger"
	uc "github.com/Manish-Mehta/tigerhall/internal/controller/user"
	"github.com/Manish-Mehta/tigerhall/internal/middleware"
	ss "github.com/Manish-Mehta/tigerhall/internal/service/sight"
	ts "github.com/Manish-Mehta/tigerhall/internal/service/tiger"
	us "github.com/Manish-Mehta/tigerhall/internal/service/user"
	"github.com/Manish-Mehta/tigerhall/model/datastore"
	"github.com/Manish-Mehta/tigerhall/pkg/db"
	"github.com/gin-gonic/gin"
)

func SetupRouter(engine *gin.Engine) {
	log.Println("Initializing Routes")
	dBClient := db.GetDBClient().GetClient()
	engine.MaxMultipartMemory = 8 << 20 // 8 MiB

	apiRouter := engine.Group("/api/v1")
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
				userRouter.POST("/refresh", middleware.AuthMiddleware, userController.Refresh)

			}
		}
		// Tiger Router
		{
			tigerRouter := apiRouter.Group("/tiger")
			tigerStore := datastore.NewTigerStore(dBClient)
			tigerService := ts.NewTigerService(tigerStore)

			tigerController := tc.NewTigerController(tigerService)
			{
				tigerRouter.POST("", middleware.AuthMiddleware, tigerController.Create)
			}
		}
		// Sight Router
		{
			sightRouter := apiRouter.Group("/sight")
			sightStore := datastore.NewSightStore(dBClient)
			sightService := ss.NewSightService(sightStore)

			sightController := sc.NewSightController(sightService)
			{
				sightRouter.POST("", middleware.AuthMiddleware, sightController.Create)
			}
		}
	}
}

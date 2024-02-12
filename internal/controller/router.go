package controller

import (
	"log"

	uc "github.com/Manish-Mehta/tigerhall/internal/controller/user"
	"github.com/Manish-Mehta/tigerhall/internal/middleware"
	us "github.com/Manish-Mehta/tigerhall/internal/service/user"
	"github.com/Manish-Mehta/tigerhall/model/datastore"
	"github.com/Manish-Mehta/tigerhall/pkg/db"
	"github.com/gin-gonic/gin"
)

func SetupRouter(engine *gin.Engine) {
	log.Println("Initializing Routes")
	dBClient := db.GetDBClient().GetClient()

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
				userRouter.POST("/", userController.Signup)
				userRouter.POST("/login", userController.Login)
				userRouter.POST("/refresh", middleware.AuthMiddleware, userController.Refresh)

			}
		}
		// Tiger Router
		{
			// userRouter := apiRouter.Group("/tiger")
			// userStore := datastore.NewUserStore(dBClient)
			// userService := us.NewUserService(userStore)

			// userController := uc.NewUserController(userService)
			// {
			// 	userRouter.POST("/", userController.Signup)
			// }
		}
	}
}

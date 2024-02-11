package controller

import (
	"log"

	uc "github.com/Manish-Mehta/tigerhall/internal/controller/user"
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

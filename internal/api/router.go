package api

import (
	"log"

	"github.com/Manish-Mehta/tigerhall/internal/api/user"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	log.Println("Initializing Routes")

	apiRouter := router.Group("/api/v1")
	{
		apiRouter.HEAD("/health", HealthApi)
		apiRouter.GET("/health", HealthApi)

		userRouter := apiRouter.Group("/user")
		{
			userRouter.POST("/", user.CreateUser())

		}
		// tigerRouter := apiRouter.Group("/tiger")
		// {

		// }
	}
}

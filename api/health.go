package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

func HealthApi(c *gin.Context) {

	returnObj := struct {
		Message string `json:"message"`
		Date    string `json:"date"`
		Uptime  string `json:"uptime"`
	}{
		Message: "Ok",
		Date:    time.Now().String(),
		Uptime:  time.Since(startTime).String(),
	}

	c.JSON(200, returnObj)
}

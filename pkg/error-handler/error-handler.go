package errorhandler

import (
	"fmt"
	"log"

	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"
)

// FATAL ERROR LOG [DO NOT USE IN API CODE]
func CheckErrorAndExit(err error, msg string) {
	if err != nil {
		log.Println(err)
		log.Fatal(msg)
	}
}

// GinDefaultRecoveryMiddelware recovers from the panic and sends the "error response" with "default error message and status code(i.e, 500)" to the client.
func GinDefaultRecoveryMiddelware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Println(fmt.Sprintf("Recovered the panic: %v", recovered))
		interceptor.SendErrRes(c, recovered, interceptor.DEFAULT_ERROR_MSG, interceptor.DEFAULT_HTTP_ERROR_CODE)
	})
}

// Checks if a panic occured and then recover along with logging it and return the recovered object.
func RecoverErr(errMsg string) interface{} {
	if r := recover(); r != nil {
		if errMsg == "" {
			errMsg = interceptor.DEFAULT_ERROR_MSG
		}
		log.Println(fmt.Sprintf(errMsg))
		return r
	}
	return nil
}

// RecoverAndSendErrRes recovers from the panic and sends the "error response" with "500 status code" to the client.
// If errMsg is not provided then "default error message" will be used.
func RecoverAndSendErrRes(c *gin.Context, errMsg string) {
	if r := RecoverErr(errMsg); r != nil {
		interceptor.SendErrRes(c, r, errMsg, interceptor.DEFAULT_HTTP_ERROR_CODE)
	}
}

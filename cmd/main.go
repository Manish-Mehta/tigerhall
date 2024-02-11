package main

import (
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/Manish-Mehta/tigerhall/internal/api"
	"github.com/Manish-Mehta/tigerhall/internal/config"
	"github.com/Manish-Mehta/tigerhall/pkg/db"
	"github.com/Manish-Mehta/tigerhall/pkg/email"
	"github.com/Manish-Mehta/tigerhall/pkg/resty"
)

func init() {
	// LoadEnv is not required in managed apps deployment where variables are stored in their key vaults(ex: Kuberenetes secrets)
	config.LoadEnv()
	config.SetConfig()
}

func main() {

	server := config.Server{GinInstance: nil}
	server.InitServer()

	resty.CreateRestyClient()

	// Initiate Brevo Email Service
	email.InitService(config.EMAIL_SERVICE)
	email.GetServiceClient(config.EMAIL_SERVICE).CreateClient(email.ClientParam{})

	// Connect DB
	db.InitService()
	// TODO: Add migration

	allowedOrigins := config.ALLOWED_ORIGINS
	if allowedOrigins == "" {
		log.Fatal("Allowed Origins not defined")
	}
	server.AddCors(&cors.Config{
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "ResponseType"},
		// ExposeHeaders:    []string{"Content-Length"},
		AllowOrigins:     strings.Split(allowedOrigins, ","),
		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	})

	api.SetupRouter(server.GinInstance)

	// must be last line
	server.Listen()
}

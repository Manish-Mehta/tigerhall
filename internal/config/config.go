package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
)

var (
	DB_STR          string
	SERVER_PORT     string
	ALLOWED_ORIGINS string

	TOKEN_SECRET string

	// Currently set to "brevo" as using brevo
	EMAIL_SERVICE      string
	EMAIL_API_ENDPOINT string
	EMAIL_FROM_ADDRESS string
	EMAIL_API_KEY      string
	// AWS Config[Not used]
	AWS_REGION            string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
)

// SHOULD BE THE FIRST FUNCTION CALL IN THE SERVICE BOOT PROCESS
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
		// FATAL LOG WILL CLOSE THE APPLICATION
		log.Fatal("Error loading ENV config file")
	}
	log.Println("ENV params loaded")
}

var Getenv = func(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("env variable %s not found", key)
	}
	return strings.Trim(val, " ")
}

func SetConfig() {
	{
		DB_STR = Getenv("DB_STR")
		SERVER_PORT = Getenv("SERVER_PORT")
		ALLOWED_ORIGINS = Getenv("ALLOWED_ORIGINS")
	}
	{
		TOKEN_SECRET = Getenv("TOKEN_SECRET")
	}
	// AWS Config[Not used as Brevo is used for email]
	{
		AWS_REGION = Getenv("AWS_REGION")
		AWS_ACCESS_KEY_ID = Getenv("AWS_ACCESS_KEY_ID")
		AWS_SECRET_ACCESS_KEY = Getenv("AWS_SECRET_ACCESS_KEY")
	}
	{
		EMAIL_API_ENDPOINT = Getenv("EMAIL_API_ENDPOINT")
		EMAIL_SERVICE = Getenv("EMAIL_SERVICE")
		EMAIL_FROM_ADDRESS = Getenv("EMAIL_FROM_ADDRESS")
		EMAIL_API_KEY = Getenv("EMAIL_API_KEY")
	}
}

type Server struct {
	GinInstance *gin.Engine
}

func (s *Server) InitServer() {
	log.Println("Initializing server")

	s.GinInstance = gin.Default()
	s.GinInstance.SetTrustedProxies(nil)
	s.GinInstance.Use(errorHandler.GinDefaultRecoveryMiddelware())
}

func (s *Server) Listen() {
	port := SERVER_PORT
	if port == "" {
		port = "3000"
	}

	s.GinInstance.Run(":" + port)
}

func (s *Server) AddCors(config *cors.Config) {

	corsConfig := config
	if corsConfig == nil {
		corsConfig = &cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "ResponseType"},
			AllowCredentials: true,

			MaxAge: 12 * time.Hour,
		}
	}

	corsObject := cors.New(*corsConfig)
	s.GinInstance.Use(corsObject)

	log.Println("CORS added")
}

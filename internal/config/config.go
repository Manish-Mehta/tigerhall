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

type EmailEvent struct {
	TigerName  string
	UserEmails []*string
}

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

	IMAGE_STORAGE_PATH string

	TIGER_SIGHTING_CHAN = make(chan uint, 10)
	EMAIL_EVENT_CHAN    = make(chan EmailEvent, 10)
)

const (
	MAX_UPLOAD_IMAGE_SIZE int64 = 1024 * 1024 * 6 // 6MB
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

func SetConfig(getEnv func(key string) string) {
	{
		DB_STR = getEnv("DB_STR")
		SERVER_PORT = getEnv("SERVER_PORT")
		ALLOWED_ORIGINS = getEnv("ALLOWED_ORIGINS")
	}
	{
		TOKEN_SECRET = getEnv("TOKEN_SECRET")
	}
	// AWS Config[Not used as Brevo is used for email]
	{
		AWS_REGION = getEnv("AWS_REGION")
		AWS_ACCESS_KEY_ID = getEnv("AWS_ACCESS_KEY_ID")
		AWS_SECRET_ACCESS_KEY = getEnv("AWS_SECRET_ACCESS_KEY")
	}
	{
		EMAIL_API_ENDPOINT = getEnv("EMAIL_API_ENDPOINT")
		EMAIL_SERVICE = getEnv("EMAIL_SERVICE")
		EMAIL_FROM_ADDRESS = getEnv("EMAIL_FROM_ADDRESS")
		EMAIL_API_KEY = getEnv("EMAIL_API_KEY")
	}
	{
		IMAGE_STORAGE_PATH = getEnv("IMAGE_STORAGE_PATH")
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

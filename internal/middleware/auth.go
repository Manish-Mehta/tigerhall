package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Manish-Mehta/tigerhall/internal/config"
	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context) {
	var claims *jwt.RegisteredClaims
	var ok bool

	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		interceptor.SendErrRes(c, "Access token missing", "Provide a valid access token", http.StatusUnauthorized)
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.TOKEN_SECRET), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		interceptor.SendErrRes(c, "Access denied", "Invalid access token", http.StatusUnauthorized)
	} else if claims, ok = token.Claims.(*jwt.RegisteredClaims); ok {
		log.Println("claims.Subject, claims.ExpiresAt")
		log.Println(claims.Subject, claims.ExpiresAt)
	} else {
		interceptor.SendErrRes(c, "Token verification failed", "Error while checking token", http.StatusUnauthorized)
	}

	if strings.HasSuffix(c.Request.URL.Path, "refresh") {
		c.Set("TokenExpiry", claims.ExpiresAt)
	}
	c.Set("Email", claims.Subject)

	c.Next()
}

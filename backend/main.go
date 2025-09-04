package main

import (
	"log"
	"os"
	"platform/backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	DefaultPort = "8080"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.RedirectTrailingSlash = false

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           86400,
	}

	// Use environment-specific CORS configuration
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		corsConfig.AllowOrigins = []string{origins}
	} else {
		// Default safe origins for development
		corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	}

	r.Use(cors.New(corsConfig))

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	log.Printf("Starting server on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
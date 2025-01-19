package main

import (
	"test/db"
	"test/models"
	"test/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	models.CreateTemporaryCatalogue()

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,           // Allow cookies to be sent with the request
		MaxAge:           12 * time.Hour, // Cache preflight response for 12 hours
	}))

	routes.RegisterRoutes(server)

	server.Run(":8080")
}

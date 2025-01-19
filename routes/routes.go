package routes

import (
	"test/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/singup", singup)
	server.POST("/login", login)
	server.GET("/catalogue", getCatalogue)
	server.POST("/guest", createGuestUser)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authentificate)

	authenticated.POST("/favorites", addToFavorites)
	authenticated.DELETE("/favorites", deleteFromFavorites)
	authenticated.GET("/favorites", getFavorites)
	authenticated.GET("/homeCatalogues", getHomeCatalogue)

}

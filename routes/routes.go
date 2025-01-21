package routes

import (
	"test/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/login", login)           //OK
	server.POST("/guest", createGuestUser) //OK

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authentificate)

	authenticated.POST("/favorite", addToFavorites)        //OK
	authenticated.DELETE("/favorite", deleteFromFavorites) //OK
	authenticated.GET("/favorite", getFavorites)           //OK

	authenticated.POST("/singup", singup) //OK

	authenticated.GET("/homeCatalogues", getHomeCatalogue) //OK
	authenticated.GET("/catalogue", getCatalogue)          //OK

	authenticated.GET("/cart", getCart)           //OK
	authenticated.POST("/cart", addToCart)        //OK
	authenticated.PUT("/cart", removeFromCart)    //OK
	authenticated.DELETE("/cart", deleteFromCart) //OK
}

package routes

import (
	"test/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/login", login)
	server.POST("/guest", createGuestUser)
	server.POST("/all", getAll)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authentificate)

	authenticated.POST("/loginJWT", loginJWT)
	authenticated.POST("/googleAuth", googleAuth)
	authenticated.POST("/userInfo", userInfo)

	authenticated.POST("/verifyEmail", verifyEmail)

	authenticated.POST("/favorite", addToFavorites)
	authenticated.DELETE("/favorite", deleteFromFavorites)
	authenticated.GET("/favorite", getFavorites)

	authenticated.GET("/orders", getOrders)

	authenticated.POST("/adress", addToAdresses)
	authenticated.DELETE("/adress", deleteFromAdresses)
	authenticated.GET("/adress", getAdresses)

	authenticated.POST("/singup", singup)

	authenticated.GET("/homeCatalogues", getHomeCatalogue)
	authenticated.GET("/catalogue", getCatalogue)

	authenticated.GET("/cart", getCart)
	authenticated.POST("/cart", addToCart)
	authenticated.PUT("/cart", removeFromCart)
	authenticated.DELETE("/cart", deleteFromCart)

	authenticated.POST("/checkout", checkoutURL)
	authenticated.POST("/checkoutSuccessfull", checkoutSuccessful)
	authenticated.POST("/checkoutFailed", checkoutFailed)

}

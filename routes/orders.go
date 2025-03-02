package routes

import (
	"net/http"
	"test/models"

	"github.com/gin-gonic/gin"
)

func checkoutURL(context *gin.Context) {

	var order = models.Order{
		UserID: context.GetInt64("userID"),
	}

	err := context.ShouldBindJSON(&order)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request data"})
		return
	}

	checkoutURL, err := order.Payment()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not get the checkout URL", "err": err.Error()})
		return
	}

	err = order.CreateOrder()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create the order", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The checkout URL was successfuly fetched", "url": checkoutURL})

}

func checkoutSuccessful(context *gin.Context) {

	var order models.Order

	err := context.ShouldBindJSON(&order)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request data"})
		return
	}

	order.UserID = context.GetInt64("userID")

	err = order.PaymentSuccessful()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not save the order", "err": err.Error()})
		return
	}

	var cart = models.UserCart{
		UserID: order.UserID,
	}

	err = cart.Clean()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not clean the cart"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The order  was successfuly saved"})

}

func checkoutFailed(context *gin.Context) {

	var order models.Order

	err := context.ShouldBindJSON(&order)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request data"})
		return
	}

	order.UserID = context.GetInt64("userID")

	err = order.PaymentFailed()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not save the order", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The order  was successfuly saved"})

}

func getOrders(context *gin.Context) {

	orders, err := models.GetOrders(context.GetInt64("userID"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not reach to the data", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"orders": orders})
}

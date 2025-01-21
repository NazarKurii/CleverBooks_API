package routes

import (
	"net/http"
	"strconv"
	"test/models"

	"github.com/gin-gonic/gin"
)

type CartRequest struct {
	BookID int64 `json:"bookId"`
}

func getCart(context *gin.Context) {

	userID := context.GetInt64("userID")

	var userCart = models.UserCart{
		UserID: userID,
	}

	err := userCart.GetUsersCart()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not get the cart"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"cart": userCart.Cart})

}

func addToCart(context *gin.Context) {

	var CartRequest CartRequest
	err := context.ShouldBindJSON(&CartRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the data"})
		return
	}

	var cart = models.CartItem{
		UserID: context.GetInt64("userID"),
		BookID: CartRequest.BookID,
	}

	err = cart.AddToCart()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not add the item to the cart"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The item was successfuly added to the cart"})

}

func removeFromCart(context *gin.Context) {
	var CartRequest CartRequest
	err := context.ShouldBindJSON(&CartRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the data"})
		return
	}
	var cart = models.CartItem{
		UserID: context.GetInt64("userID"),
		BookID: CartRequest.BookID,
	}

	err = cart.RemoveFromCart()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not remove the item from the cart"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The item was successfuly removed from the cart"})

}

func deleteFromCart(context *gin.Context) {

	bookIDString := context.Query("bookId")

	if bookIDString == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing 'bookId' parameter in the request",
		})
		return
	}

	bookID, err := strconv.Atoi(bookIDString)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "'bookId' must be a valid integer",
		})
		return
	}

	var cart = models.CartItem{
		UserID: context.GetInt64("userID"),
		BookID: int64(bookID),
	}

	err = cart.DeleteFromCart()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not delete the item from the cart"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The item was successfuly deleted from the cart"})

}

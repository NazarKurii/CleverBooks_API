package routes

import (
	"net/http"
	"strconv"
	"test/models"

	"github.com/gin-gonic/gin"
)

type favoriteRequest struct {
	BookID int64 `json:"bookId"`
}

func addToFavorites(context *gin.Context) {

	var favoriteRequest favoriteRequest
	err := context.ShouldBindJSON(&favoriteRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the data"})
		return
	}

	var favorite = models.Favorite{
		UserID:     context.GetInt64("userID"),
		FavoriteID: favoriteRequest.BookID,
	}

	err = favorite.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not add the data to favorites"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The data was saved to favorites successfuly"})

}

func deleteFromFavorites(context *gin.Context) {
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

	var favorite = models.Favorite{
		UserID:     context.GetInt64("userID"),
		FavoriteID: int64(bookID),
	}

	err = favorite.Delete()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not delete the data from favorites"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The data was deleted from favorites successfuly"})

}

func getFavorites(context *gin.Context) {

	favorites, err := models.GetFavorites(context.GetInt64("userID"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not reach to the data", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"books": favorites})
}

package routes

import (
	"net/http"
	"test/models"

	"github.com/gin-gonic/gin"
)

type favoriteRequest struct {
	FavoriteID int64 `json:"favorite_id"`
}

func addToFavorites(context *gin.Context) {

	var favoriteRequest favoriteRequest
	err := context.ShouldBindJSON(&favoriteRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the data", "error": err.Error()})
		return
	}

	var favorite = models.Favorite{
		UserID:     context.GetInt64("userID"),
		FavoriteID: favoriteRequest.FavoriteID,
	}

	err = favorite.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not add the data to favorites"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The data was saved to favorites successfuly"})

}

func deleteFromFavorites(context *gin.Context) {

	var favoriteRequest favoriteRequest
	err := context.ShouldBindJSON(&favoriteRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the data", "error": err.Error()})
		return
	}

	var favorite = models.Favorite{
		UserID:     context.GetInt64("userID"),
		FavoriteID: favoriteRequest.FavoriteID,
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
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not reach to the data"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

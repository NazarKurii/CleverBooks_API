package routes

import (
	"net/http"
	"strconv"
	"test/models"

	"github.com/gin-gonic/gin"
)

type AdressRequest struct {
	ID int64 `json:"id"`
}

func addToAdresses(context *gin.Context) {

	var adress models.Adress
	err := context.ShouldBindJSON(&adress)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the data"})
		return
	}

	adress.UserID = context.GetInt64("userID")

	id, err := adress.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not add the data to adresses", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The data was saved to adresses successfuly", "id": id})

}

func deleteFromAdresses(context *gin.Context) {
	adressIDString := context.Query("id")

	if adressIDString == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing 'id' parameter in the request",
		})
		return
	}

	adressID, err := strconv.Atoi(adressIDString)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "'id' must be a valid integer",
		})
		return
	}

	adress := models.Adress{
		ID:     int64(adressID),
		UserID: context.GetInt64("userID"),
	}

	err = adress.Delete()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not delete the data from adresses", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "The data was deleted from adresses successfuly"})

}

func getAdresses(context *gin.Context) {

	adresses, err := models.GetAdresses(context.GetInt64("userID"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not reach to the data", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"adresses": adresses})
}

package routes

import (
	"net/http"
	"test/models"
	"test/utils"

	"github.com/gin-gonic/gin"
)

func singup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request data"})
		return
	}

	token := context.Request.Header.Get("Authorization")

	err = user.Save(token)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not save the user"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "The user was created successfully"})

}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request data"})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateUserToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request data"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func createGuestUser(context *gin.Context) {

	var user models.User

	err := user.CreateGuest()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create a guest user"})
		return
	}

	token, err := utils.GenerateGuestToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create a guest user"})
		return
	}

	err = user.Save("")

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create a guest user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Guesr user was successfuly created", "token": token})
}

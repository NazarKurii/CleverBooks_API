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

	user.ID = context.GetInt64("userID")

	token, err := user.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not save the user", "err": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "The user was saved successfuly", "token": token})

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

func loginJWT(context *gin.Context) {

	user := models.User{
		Email: context.GetString("email"),
		ID:    context.GetInt64("userID"),
	}

	err := user.IsRegistered()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not log the user in", "err": err.Error()})
		return
	}

	var token string
	var message string

	if user.Registered {
		token, err = utils.GenerateUserToken(user.Email, user.ID)
		message = "User was successfuly logged in"
	} else {
		token, err = utils.GenerateGuestToken(user.Email, user.ID)
		message = "Guest user was successfuly logged in"
	}

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not generate a token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": message, "token": token, "registered": user.Registered})

}

func createGuestUser(context *gin.Context) {

	var user models.User

	token, err := user.CreateGuest()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create a guest user 1"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Guest user was successfuly created", "token": token})

}

func verifyEmail(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request data"})
		return
	}

	emailExists, err := user.VerifyEmail()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not verify the email", "err": err.Error()})
		return
	}

	var verificationCode string

	if !emailExists {

		verificationCode, err = utils.SendVerificationCode(user.Email)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Could not verify the email", "err": err.Error()})
			return
		}
	}

	context.JSON(http.StatusOK, gin.H{"message": "Number was successfuly verified", "exists": emailExists, "verificationCode": verificationCode})

}

func getUser(context *gin.Context) {

	user := models.User{
		Email: context.GetString("email"),
	}

	err := user.GetUser()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found", "err": err.Error()})
		return
	}

	context.JSON(http.StatusBadRequest, gin.H{"message": "User found", "user": user})
}

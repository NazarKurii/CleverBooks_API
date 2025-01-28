package middlewares

import (
	"net/http"
	"test/utils"

	"github.com/gin-gonic/gin"
)

func Authentificate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	userID, email, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	context.Set("userID", userID)
	context.Set("email", email)

	context.Next()

}

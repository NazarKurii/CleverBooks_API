package routes

import (
	"net/http"
	"test/db"
	"test/models"
	"test/utils"

	"github.com/gin-gonic/gin"
)

func getHomeCatalogue(context *gin.Context) {

	rows, err := db.DB.Query("SELECT id FROM catalogue")

	if err != nil {
		context.JSON(http.StatusOK, gin.H{"message": "Could not get the catalogue "})
		return
	}

	var IDs []int

	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			context.JSON(http.StatusOK, gin.H{"message": "Could not get the catalogue "})
			return
		}

		IDs = append(IDs, id)
	}

	var catalogue models.Catalogue

	err = catalogue.GetBooksInfo(IDs, context.GetInt64("userID"))

	if err != nil {
		context.JSON(http.StatusOK, gin.H{"message": "Could not get the catalogue", "error": err.Error()})
		return
	}

	catalogues := models.Catalogue(catalogue).Sort()

	context.JSON(http.StatusOK, gin.H{"sections": catalogues})
}

func getCatalogue(context *gin.Context) {

	param := context.QueryArray("ids")

	IDs, err := utils.StringToInt(param)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data", "err": err.Error()})
		return
	}

	var catalogue models.Catalogue

	err = catalogue.GetBooksInfo(IDs, context.GetInt64("userID"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not find books", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"books": catalogue})

}

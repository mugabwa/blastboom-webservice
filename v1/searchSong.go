package v1

import (
	"net/http"

	api "blastboom/webservice/apis"

	"github.com/gin-gonic/gin"
)

func SearchHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		query := ctx.Query("q")
		if query == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
			return
		}
		accToken, _ := tokenMx.GetToken()
		

		results, err := api.SearchSpotify(accToken, query, "track", 10)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, results)
	}
}
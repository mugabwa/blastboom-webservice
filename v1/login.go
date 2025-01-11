package v1

import (
	api "blastboom/webservice/apis"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func UserLogin(ctx *gin.Context) {
	params := url.Values{}
	params.Set("client_id", api.ClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", api.RedirectURI)
	// user-modify-playback-state streaming
	params.Set("scope", "user-read-email user-read-private user-read-playback-state user-modify-playback-state")
	authURL := fmt.Sprintf("%s?%s", api.BaseAuthURL, params.Encode())
	ctx.Redirect(http.StatusFound, authURL)
}

func HandleCallback(tokenMx *api.TokenManager) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		code := ctx.Query("code")
		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing authorization code"})
			return
		}
		token, err := api.ExchangeAccessToken(code)
		tokenMx.SetToken(token.AccessToken, int(token.ExpiresIn))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		accToken, _ := tokenMx.GetToken()
		fmt.Print(accToken)
		profile, err := api.GetProfile(accToken, api.BaseAPIURL)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, profile)
	}
}
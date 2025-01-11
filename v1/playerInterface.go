package v1

import (
	api "blastboom/webservice/apis"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PlayBackHandler handles the playback retrieval by validating the token and calling the GetPlayBack API.
// It takes a TokenManager as a parameter to manage the access token.
// If the token is invalid, it returns a 401 Unauthorized response.
// If the token is valid, it calls the GetPlayBack API and returns the results or an error.
func PlayBackHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		results, httpStatus, err := api.GetPlayBack(accessToken)
		if err != nil {
			ctx.JSON(httpStatus, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, results)
	}
}

// PlayBackTransferHandler handles the transfer of playback to a different device by validating the token and calling the TransferPlayback API.
// Parameters:
// - tokenMx: a pointer to the TokenManager instance used to manage access tokens.
func PlayBackTransferHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		deviceID := ctx.PostForm("device_id")
		play := ctx.PostForm("play") == "true"
		if deviceID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
			return
		}
		status, err := api.TransferPlayback(accessToken, deviceID, play)
		if err != nil {
			ctx.JSON(status, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(status, gin.H{"status": "Playback transferred"})
	}
}

// DevicesHandler handles the retrieval of available devices by validating the token and calling the GetDevices API.
// It takes a TokenManager as a parameter to manage the access token.
// If the token is invalid, it returns a 401 Unauthorized response.
// If the token is valid, it calls the GetDevices API and returns the results or an error.
func DevicesHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		result, status, err := api.GetDevices(accessToken)
		if err != nil {
			ctx.JSON(status, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}

func CurrentPlayingTrackHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		result, status, err := api.GetCurrentPlayingTrack(accessToken)
		if err != nil {
			ctx.JSON(status, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(status, result)
	}
}

func StartPlaybackHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		var json struct {
			DeviceID string `json:"device_id,omitempty"`
			ContextURI string `json:"context_uri"`
			Offset *struct {
				Position int `json:"position,omitempty"`
			} `json:"offset,omitempty"`
			PositionMS int `json:"position_ms,omitempty"`
		}
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		deviceID := json.DeviceID
		contextURI := json.ContextURI
		offsetPosition := 0
		if json.Offset != nil {
			offsetPosition = json.Offset.Position
		}
		positionMS := json.PositionMS

		statusCode, err := api.StartPlayback(
			accessToken, deviceID, contextURI, offsetPosition, positionMS)
		if err != nil {
			ctx.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(statusCode, gin.H{"status": "Playback started"})
	}
}

func PausePlaybackHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		var json struct {
			DeviceID string `json:"device_id"`
		}
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		deviceID := json.DeviceID
		statusCode, err := api.PausePlayback(accessToken, deviceID)
		if err != nil {
			ctx.JSON(statusCode, gin.H{"error": err.Error()})
		}
		ctx.JSON(statusCode, gin.H{"status": "Playback paused"})
	}
}

func SkipNextHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		var json struct {
			DeviceID string `json:"device_id"`
		}
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		deviceID := json.DeviceID
		statusCode, err := api.SkipNext(accessToken, deviceID)
		if err != nil {
			ctx.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(statusCode, gin.H{"status": "Skipped to next track"})
	}
}

func SkipPrevHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		var json struct {
			DeviceID string `json:"device_id"`
		}
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		deviceID := json.DeviceID
		statusCode, err := api.SkipPrev(accessToken, deviceID)
		if err != nil {
			ctx.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(statusCode, gin.H{"status": "Skipped to previous track"})
	}
}

func SeekPositionHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		var json struct {
			DeviceID string `json:"device_id"`
			PositionMS int `json:"position_ms"`
		}
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		deviceID := json.DeviceID
		positionMS := json.PositionMS
		statusCode, err := api.SeekPosition(accessToken, deviceID, positionMS)
		if err != nil {
			ctx.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(statusCode, gin.H{"status": "Position seeked"})
	}
}

func ToggleRepeatHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		var json struct {
			DeviceID string `json:"device_id"`
			State    string `json:"state"` // "track", "context", or "off"
		}
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		deviceID := json.DeviceID
		state := json.State
		statusCode, err := api.ToggleRepeat(accessToken, deviceID, state)
		if err != nil {
			ctx.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(statusCode, gin.H{"status": "Repeat state toggled"})
	}
}

func SetPlaybackVolumeHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		var json struct {
			DeviceID string `json:"device_id"`
			Volume   int    `json:"volume"`
		}
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		deviceID := json.DeviceID
		volume := json.Volume
		statusCode, err := api.SetPlaybackVolume(accessToken, deviceID, volume)
		if err != nil {
			ctx.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(statusCode, gin.H{"status": "Playback volume set"})
	}
}

func ToggleShuffleHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		var json struct {
			DeviceID string `json:"device_id"`
			State    bool   `json:"state"` // true for shuffle on, false for shuffle off
		}
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		deviceID := json.DeviceID
		state := json.State
		statusCode, err := api.ToggleShuffle(accessToken, deviceID, state)
		if err != nil {
			ctx.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(statusCode, gin.H{"status": "Shuffle state toggled"})
	}
}

func GetRecentlyPlayedHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		results, status, err := api.GetRecentlyPlayed(accessToken)
		if err != nil {
			ctx.JSON(status, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, results)
	}
}

func GetUsersQueueHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		results, status, err := api.GetUsersQueue(accessToken)
		if err != nil {
			ctx.JSON(status, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, results)
	}
}

func AddToQueueHandler(tokenMx *api.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, valid := tokenMx.GetToken()
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		var json struct {
			DeviceID string `json:"device_id"`
			URI      string `json:"uri"`
		}
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		deviceID := json.DeviceID
		uri := json.URI
		statusCode, err := api.AddToQueue(accessToken, deviceID, uri)
		if err != nil {
			ctx.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(statusCode, gin.H{"status": "Added to queue"})
	}
}

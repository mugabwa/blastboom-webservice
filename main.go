package main

import (
	"github.com/gin-gonic/gin"

	"blastboom/webservice/v1"
	"blastboom/webservice/apis"
)


func main() {
	tokenManager := api.NewTokenManager()
	router := gin.Default()
	router.GET("/login", v1.UserLogin)
	router.GET("/callback", v1.HandleCallback(tokenManager))
	router.GET("/search", v1.SearchHandler(tokenManager))
	router.GET("/player", v1.PlayBackHandler(tokenManager))
	router.PUT("/player", v1.PlayBackTransferHandler(tokenManager))
	router.GET("/player/devices", v1.DevicesHandler(tokenManager))
	router.GET("/player/currently-playing", v1.CurrentPlayingTrackHandler(tokenManager))
	router.PUT("/player/play", v1.StartPlaybackHandler(tokenManager))
	router.PUT("/player/pause", v1.PausePlaybackHandler(tokenManager))
	router.PUT("/player/next", v1.SkipNextHandler(tokenManager))
	router.PUT("/player/previous", v1.SkipPrevHandler(tokenManager))
	router.PUT("/player/seek", v1.SeekPositionHandler(tokenManager))
	router.PUT("/player/repeat", v1.ToggleRepeatHandler(tokenManager))
	router.PUT("/player/volume", v1.SetPlaybackVolumeHandler(tokenManager))
	router.PUT("/player/shuffle", v1.ToggleShuffleHandler(tokenManager))
	router.GET("/player/recently-played", v1.GetRecentlyPlayedHandler(tokenManager))
	router.GET("/player/queue", v1.GetUsersQueueHandler(tokenManager))
	router.POST("/player/queue", v1.AddToQueueHandler(tokenManager))
	router.Run("localhost:8080")
}


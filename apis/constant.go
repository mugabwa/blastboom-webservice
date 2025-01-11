package api

import (
	"os"
)

var (
	ClientID     = os.Getenv("SPOTIFY_CLIENT_ID")
	ClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	RedirectURI  = os.Getenv("SPOTIFY_REDIRECT_URI")
)

const (
	BaseAuthURL  = "https://accounts.spotify.com/authorize"
	BaseTokenURL = "https://accounts.spotify.com/api/token"
	BaseAPIURL   = "https://api.spotify.com/v1"
)

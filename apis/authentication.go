package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type accessPayload struct {
	AccessToken string `json:"access_token"`
	TokenType 	string `json:"token_type"`
	ExpiresIn 	int32 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func ExchangeAccessToken(code string) (*accessPayload, error) {
	authURL := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", ClientID)
	data.Set("client_secret", ClientSecret)
	data.Set("redirect_uri", RedirectURI)

	req, err := http.NewRequest("POST", authURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("authentication failed: %s", body)
	}

	var payload accessPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &payload, nil
}

type UserProfile struct {
	DisplayName string `json:"display_name"`
	ID			string `json:"id"`
	Email		string `json:"email"`
	Country		string `json:"country"`
	Product		string `json:"product"`
	Uri 		string `json:"uri"`
}
 
func GetProfile(accessToken, baseUrl string) (*UserProfile, error) {
	url := baseUrl + "/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + accessToken)
	// req.Header.Set("Content-Type", "application/json")
	fmt.Println(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch profile: %s", body)
	}

	var profile UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &profile, nil
}

type TokenManager struct {
	AccessToken string
	ExpiresAt time.Time
	mutx	sync.RWMutex
}

func NewTokenManager() *TokenManager {
	return &TokenManager{}
}

func (tm *TokenManager) SetToken(token string, expiresIn int) {
	tm.mutx.Lock()
	defer tm.mutx.Unlock()

	tm.AccessToken = token
	tm.ExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

func (tm *TokenManager) GetToken() (string, bool) {
	tm.mutx.RLock()
	defer tm.mutx.RUnlock()

	if time.Now().Before(tm.ExpiresAt) {
		return tm.AccessToken, true
	}
	return "", false
}

func (tm *TokenManager) RefreshToken(refreshFunc func() (string, int, error)) (string, error) {
	tm.mutx.Lock()
	defer tm.mutx.Unlock()

	if time.Now().Before(tm.ExpiresAt) {
		return tm.AccessToken, nil
	}

	token, expiresIn, err := refreshFunc()
	if err != nil {
		return "", err
	}

	tm.AccessToken = token
	tm.ExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return token, nil
}

func (tm *TokenManager)GetAccessToken(code string) (string, error) {
	tm.mutx.Lock()
	defer tm.mutx.Unlock()

	fmt.Println(tm)
	if time.Now().Before(tm.ExpiresAt) {
		return tm.AccessToken, nil
	}

	token, err := ExchangeAccessToken(code)
	if err != nil {
		return "", fmt.Errorf("failed to fetch access token: %w", err)
	}

	fmt.Println(token.ExpiresIn)
	tm.AccessToken = token.AccessToken
	tm.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	fmt.Println(tm.ExpiresAt)

	return tm.AccessToken, nil
}

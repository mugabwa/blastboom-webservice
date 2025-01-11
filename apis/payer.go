package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Actions struct {
	InterruptingPlayback bool `json:"interrupting_playback"`
	Pausing bool `json:"pausing"`
	Resuming bool `json:"resuming"`
	Seeking bool `json:"seeking"`
	SkippingNext bool `json:"skipping_next"`
	SkippingPrev bool `json:"skipping_prev"`
	TogglingRepeatContext bool `json:"toggling_repeat_context"`
	TogglingShuffle bool `json:"toggling_shuffle"`
	TogglingRepeatTrack bool `json:"toggling_repeat_track"`
	TransferringPlayback bool `json:"transferring_playback"`
}

type PlayBackResponse struct {
	Device *DeviceData `json:"device"`
	RepeatState string `json:"repeat_state"`
	ShuffleState bool `json:"shuffle_state"`
	Timestamp uint64 `json:"timestamp"`
	ProgressMS uint64 `json:"progress_ms"`
	Item *Track `json:"item"`
	CurrentlyPlayingType string `json:"currently_playing_type"`
	Actions *Actions `json:"actions"`
}

func GetPlayBack(accessToken string) (*PlayBackResponse, int, error) {
	playbackURL := fmt.Sprintf("%s/me/player", BaseAPIURL)
	req, err := http.NewRequest("GET", playbackURL, nil)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, resp.StatusCode, fmt.Errorf("no content found: playback state unavailable")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, resp.StatusCode, fmt.Errorf("failed to get player: %s", body)
	}

	var results PlayBackResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to parse response: %w", err)
	}

	return &results, resp.StatusCode, nil
}

func TransferPlayback(accessToken, deviceID string, play bool) (int, error){
	transferURL := fmt.Sprintf("%s/me/player", BaseAPIURL)
	payload := map[string]interface{}{
		"device_ids": []string{deviceID},
		"play": play,
	}
	statusCode := 500
	body, err := json.Marshal(payload)
	if err != nil {
		return statusCode, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("PUT", transferURL, bytes.NewBuffer(body))
	if err != nil {
		return statusCode, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return statusCode, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("failed to transfer playback: %s", body)
	}

	return resp.StatusCode, nil
}

type DeviceData struct {
	ID string `json:"id"`
	IsActive bool `json:"is_active"`
	IsPrivateSession bool `json:"is_private_session"`
	IsRestricted bool `json:"is_restricted"`
	Name string `json:"name"`
	Type string `json:"type"`
	VolumePercent int32 `json:"volume_percent"`
	SupportsVolume bool `json:"supports_volume"`
}

type DeviceResponse struct {
	Devices []DeviceData `json:"devices"`
}

func GetDevices(accessToken string) (*DeviceResponse, int, error) {
	deviceURL := fmt.Sprintf("%s/me/player/devices", BaseAPIURL)
	status := 500
	req, err := http.NewRequest("GET", deviceURL, nil)
	if err != nil {
		return nil, status, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, status, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, resp.StatusCode, fmt.Errorf("failed to get devices: %s", body)
	}

	var results DeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to parse response: %w", err)
	}

	return &results, resp.StatusCode, nil
}

type CurrentTrackResponse struct {
	Device *DeviceData `json:"device"`
	RepeatState string `json:"repeat_state"`
	ShuffleState bool `json:"shuffle_state"`
	ProgressMS int64 `json:"progress_ms"`
	IsPlaying bool `json:"is_playing"`
	Item *Track `json:"item"`
	Actions *Actions `json:"actions"`
}

func GetCurrentPlayingTrack(accessToken string) (*CurrentTrackResponse, int, error) {
	trackURL := fmt.Sprintf("%s/me/player/currently-playing", BaseAPIURL)
	statusCode := 500
	req, err := http.NewRequest("GET", trackURL, nil)
	if err != nil {
		return nil, statusCode, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, statusCode, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, resp.StatusCode, fmt.Errorf("failed to get track: %s", body)
	}

	var results CurrentTrackResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to parse response: %w", err)
	}

	return &results, resp.StatusCode, nil
}

func StartPlayback(accessToken, deviceID, contextURI string, offsetPosition, positionMS int) (int, error) {
	trackURL := fmt.Sprintf("%s/me/player/play?device_id=%s", BaseAPIURL, deviceID)
	statusCode := 500
	payload := map[string]interface{}{
		"context_uri": contextURI, 
		"offset": map[string]interface{}{
			"position": offsetPosition,
		},
		"position_ms": positionMS,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return statusCode, fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err := http.NewRequest("PUT", trackURL, bytes.NewBuffer(body))
	if err != nil {
		return statusCode, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return statusCode, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if (resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK) {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("failed to start playback: %s", body)
	}

	return resp.StatusCode, nil
}

func PausePlayback(accessToken, deviceID string) (int, error) {
	trackURL := fmt.Sprintf("%s/me/player/pause", BaseAPIURL)
	statusCode := 500
	if deviceID != "" {
		trackURL += "?device_id=" + deviceID
	}
	req, err := http.NewRequest("PUT", trackURL, nil)
	if err != nil {
		return statusCode, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return statusCode, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if (resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK) {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("failed to pause playback: %s", body)
	}

	return resp.StatusCode, nil
}

func SkipNext(accessToken, deviceID string) (int, error) {
	trackURL := fmt.Sprintf("%s/me/player/next", BaseAPIURL)
	statusCode := 500
	if deviceID != "" {
		trackURL += "?device_id=" + deviceID
	}
	req, err := http.NewRequest("POST", trackURL, nil)
	if err != nil {
		return statusCode, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return statusCode, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if (resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK) {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("failed to skip to the next track: %s", body)
	}

	return resp.StatusCode, nil
}

func SkipPrev(accessToken, deviceID string) (int, error) {
	trackURL := fmt.Sprintf("%s/me/player/previous", BaseAPIURL)
	
	statusCode := 500
	if deviceID != "" {
		trackURL += "?device_id=" + deviceID
	}
	req, err := http.NewRequest("POST", trackURL, nil)
	if err != nil {
		return statusCode, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return statusCode, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("failed to skip to the previous track: %s", body)
	}

	return resp.StatusCode, nil

}
func SeekPosition(accessToken, deviceID string, positionMS int) (int, error) {
	seekURL := fmt.Sprintf("%s/me/player/seek?position_ms=%d", BaseAPIURL, positionMS)
	if deviceID != "" {
		seekURL += "&device_id=" + deviceID
	}
	req, err := http.NewRequest("PUT", seekURL, nil)
	if err != nil {
		return 500, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("failed to seek position: %s", body)
	}

	return resp.StatusCode, nil
}

func ToggleRepeat(accessToken, deviceID string, state string) (int, error) {
	repeatURL := fmt.Sprintf("%s/me/player/repeat?state=%s", BaseAPIURL, state)
	if deviceID != "" {
		repeatURL += "&device_id=" + deviceID
	}
	req, err := http.NewRequest("PUT", repeatURL, nil)
	if err != nil {
		return 500, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("failed to toggle repeat: %s", body)
	}

	return resp.StatusCode, nil
}

func SetPlaybackVolume(accessToken, deviceID string, volumePercent int) (int, error) {
	volumeURL := fmt.Sprintf("%s/me/player/volume?volume_percent=%d", BaseAPIURL, volumePercent)
	if deviceID != "" {
		volumeURL += "&device_id=" + deviceID
	}
	req, err := http.NewRequest("PUT", volumeURL, nil)
	if err != nil {
		return 500, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("failed to set volume: %s", body)
	}

	return resp.StatusCode, nil
}

func ToggleShuffle(accessToken, deviceID string, state bool) (int, error) {
	shuffleURL := fmt.Sprintf("%s/me/player/shuffle?state=%t", BaseAPIURL, state)
	if deviceID != "" {
		shuffleURL += "&device_id=" + deviceID
	}
	req, err := http.NewRequest("PUT", shuffleURL, nil)
	if err != nil {
		return 500, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("failed to toggle shuffle: %s", body)
	}

	return resp.StatusCode, nil
}

type RecentlyPlayedResponse struct {
	Items []RecentlyPlayedItem `json:"items"`
}

type RecentlyPlayedItem struct {
	Track      *Track      `json:"track"`
	PlayedAt   string      `json:"played_at"`
	Context    *Context    `json:"context"`
}

type Context struct {
	ExternalUrls map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
}

func GetRecentlyPlayed(accessToken string) (*RecentlyPlayedResponse, int, error) {
	recentlyPlayedURL := fmt.Sprintf("%s/me/player/recently-played", BaseAPIURL)
	req, err := http.NewRequest("GET", recentlyPlayedURL, nil)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, resp.StatusCode, fmt.Errorf("failed to get recently played tracks: %s", body)
	}

	var results RecentlyPlayedResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to parse response: %w", err)
	}

	return &results, resp.StatusCode, nil
}

type QueueResponse struct {
	CurrentlyPlaying *Track   `json:"currently_playing"`
	Queue            []*Track `json:"queue"`
}

func GetUsersQueue(accessToken string) (*QueueResponse, int, error) {
	queueURL := fmt.Sprintf("%s/me/player/queue", BaseAPIURL)
	req, err := http.NewRequest("GET", queueURL, nil)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, resp.StatusCode, fmt.Errorf("failed to get user's queue: %s", body)
	}

	var results QueueResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to parse response: %w", err)
	}

	return &results, resp.StatusCode, nil
}

func AddToQueue(accessToken, deviceID, uri string) (int, error) {
	queueURL := fmt.Sprintf("%s/me/player/queue?uri=%s", BaseAPIURL, uri)
	if deviceID != "" {
		queueURL += "&device_id=" + deviceID
	}
	req, err := http.NewRequest("POST", queueURL, nil)
	if err != nil {
		return 500, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("failed to add to playback queue: %s", body)
	}

	return http.StatusCreated, nil
}

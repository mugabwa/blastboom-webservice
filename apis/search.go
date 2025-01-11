package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Track struct {
	Name string `json:"name"`
	DurationMS int32 `json:"duration_ms"`
	ID string `json:"id"`
	DataType string `json:"type"`
	Album struct {
		Name string `json:"name"`
		ID string `json:"id"`
	} `json:"album"`
	Artists []struct {
		Name string `json:"name"`
		ID string `json:"id"`
	} `json:"artists"`
}

type SearchResponse struct {
	Tracks struct {
		Items []Track `json:"items"`
	} `json:"tracks"`
}

func SearchSpotify(accessToken, query, searchType string, limit int32) (*SearchResponse, error) {
	searchURL := BaseAPIURL + "/search"
	params := url.Values{}
	params.Set("q", query)
	params.Set("type", searchType)
	params.Set("limit", fmt.Sprintf("%d", limit))

	req, err := http.NewRequest("GET", searchURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to search music: %s", body)
	}

	var results SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &results, nil
}
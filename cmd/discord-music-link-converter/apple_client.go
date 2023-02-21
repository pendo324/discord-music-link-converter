package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	applemusic "github.com/minchao/go-apple-music"
)

type AppleClient interface {
	Search(term string, types []string) (*applemusic.Search, error)
	GetAlbumById(id string) (*applemusic.Album, error)
}

type appleClient struct {
	token string
}

var _ AppleClient = (*appleClient)(nil)

func NewAppleClient() (*appleClient, error) {
	url := "https://music.apple.com"

	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %w", url, err)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response %s: %w", url, err)
	}

	re := regexp.MustCompile(`[^"]*index\..*\.js`)
	path := re.FindString(string(b))

	r, err = http.Get(url + path)
	if err != nil {
		return nil, fmt.Errorf("failed to read response %s: %w", url+path, err)
	}

	b, err = io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response %s: %w", url+path, err)
	}

	re = regexp.MustCompile(`([^"]*)"[^"]*"x-apple-jingle-correlation-key`)
	token := re.FindStringSubmatch(string(b))[1]

	return &appleClient{
		token: token,
	}, nil
}

func (c appleClient) Search(term string, types []string) (*applemusic.Search, error) {
	url := "https://api.music.apple.com/v1/catalog/us/search"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", url, err)
	}

	q := req.URL.Query()
	q.Add("term", term)

	if len(types) > 0 {
		q.Add("types", strings.Join(types, ","))
	}

	req.URL.RawQuery = q.Encode()

	req.Header.Set("origin", "https://music.apple.com")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.token))

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to %s: %w", url, err)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body of request to %s: %w", url, err)
	}

	var searchResults applemusic.Search
	err = json.Unmarshal(b, &searchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return &searchResults, nil
}

func (c appleClient) GetAlbumById(id string) (*applemusic.Album, error) {
	url := fmt.Sprintf("https://api.music.apple.com/v1/catalog/us/albums/%s", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", url, err)
	}

	req.Header.Set("origin", "https://music.apple.com")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.token))

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to %s: %w", url, err)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body of request to %s: %w", url, err)
	}

	var albums applemusic.Albums
	err = json.Unmarshal(b, &albums)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return &albums.Data[0], nil
}

func (c appleClient) GetSongById(id string) (*applemusic.Song, error) {
	url := fmt.Sprintf("https://api.music.apple.com/v1/catalog/us/songs/%s", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", url, err)
	}

	req.Header.Set("origin", "https://music.apple.com")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.token))

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to %s: %w", url, err)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body of request to %s: %w", url, err)
	}

	var songs applemusic.Songs
	err = json.Unmarshal(b, &songs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return &songs.Data[0], nil
}

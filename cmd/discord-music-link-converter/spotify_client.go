package main

import (
	"context"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type spotifyClient struct {
	clientId string
	ctx      context.Context
	secret   string
	client   *spotify.Client
}

func NewSpotifyClient(clientId string, secret string) (*spotifyClient, error) {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: secret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		return nil, err
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	return &spotifyClient{
		clientId: clientId,
		ctx:      ctx,
		secret:   secret,
		client:   client,
	}, nil
}

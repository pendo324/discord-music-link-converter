package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify/v2"
)

type spotify_album struct {
	spotify *spotifyClient
}

var _ Player = (*spotify_album)(nil)

func NewSpotifyAlbum(client *spotifyClient) *spotify_album {
	return &spotify_album{
		spotify: client,
	}
}

func (a spotify_album) Search(name string, artist string, thingType ThingType) (*ThingInfo, error) {
	term := fmt.Sprintf("%s %s", name, artist)
	res, err := a.spotify.client.Search(a.spotify.ctx, term, spotify.SearchTypeAlbum)

	if err != nil {
		return nil, fmt.Errorf("failed to search spotify for albums with term (%s): %w", term, err)
	}

	if len(res.Albums.Albums) > 0 {
		return &ThingInfo{
			Link:   fmt.Sprintf("https://open.spotify.com/album/%s", res.Albums.Albums[0].ID.String()),
			Type:   a.HandlerType(),
			Artist: res.Albums.Albums[0].Artists[0].Name,
			Name:   res.Albums.Albums[0].Name,
		}, nil
	}

	return nil, fmt.Errorf("no results")
}

func (a spotify_album) Handler(message *discordgo.MessageCreate, matches []string) *ThingInfo {
	id := a.Pattern().SubexpIndex("id")
	trackId := matches[id]

	res, err := a.spotify.client.GetTrack(a.spotify.ctx, spotify.ID(trackId))
	if err != nil {
		log.Println(fmt.Errorf("error getting the Spotify album: %w", err))
	}

	return &ThingInfo{
		Artist: res.Artists[0].Name,
		Name:   res.Name,
		Type:   a.HandlerType(),
		Link:   message.Content,
	}
}

func (spotify_album) HandlerType() ThingType {
	return ThingType("album")
}

func (spotify_album) Name() string {
	return "Spotify (album)"
}

func (spotify_album) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`https://open\.spotify\.com/album/(?P<id>[a-zA-Z0-9]+)`)
}

package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify/v2"
)

type spotify_track struct {
	spotify *spotifyClient
}

var _ Player = (*spotify_track)(nil)

func NewSpotifyTrack(client *spotifyClient) *spotify_track {
	return &spotify_track{
		spotify: client,
	}
}

func (t spotify_track) Search(name string, artist string, thingType ThingType) (*ThingInfo, error) {
	term := fmt.Sprintf("%s %s", name, artist)
	res, err := t.spotify.client.Search(t.spotify.ctx, term, spotify.SearchTypeTrack)

	if err != nil {
		return nil, fmt.Errorf("failed to search spotify for tracks with term (%s): %w", term, err)
	}

	if len(res.Tracks.Tracks) > 0 {
		return &ThingInfo{
			Link:   fmt.Sprintf("https://open.spotify.com/track/%s", res.Tracks.Tracks[0].ID.String()),
			Type:   t.HandlerType(),
			Artist: res.Tracks.Tracks[0].Artists[0].Name,
			Name:   res.Tracks.Tracks[0].Name,
		}, nil
	}

	return nil, fmt.Errorf("no results")
}

func (t spotify_track) Handler(message *discordgo.MessageCreate, matches []string) *ThingInfo {
	id := t.Pattern().SubexpIndex("id")
	trackId := matches[id]

	res, err := t.spotify.client.GetTrack(t.spotify.ctx, spotify.ID(trackId))
	if err != nil {
		log.Println(fmt.Errorf("error getting the Spotify track: %w", err))
	}

	return &ThingInfo{
		Artist: res.Artists[0].Name,
		Name:   res.Name,
		Type:   t.HandlerType(),
		Link:   message.Content,
	}
}

func (spotify_track) HandlerType() ThingType {
	return ThingType("track")
}

func (spotify_track) Name() string {
	return "Spotify (track)"
}

func (spotify_track) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`https://open\.spotify\.com/track/(?P<id>[a-zA-Z0-9]+)`)
}

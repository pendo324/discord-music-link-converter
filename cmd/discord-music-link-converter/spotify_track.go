package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify/v2"
)

type track struct {
	spotify *spotifyClient
}

var _ Player = (*track)(nil)

func NewSpotifyTrack(client *spotifyClient) *track {
	return &track{
		spotify: client,
	}
}

func (t track) Search(name string, artist string, thingType ThingType) (*ThingInfo, error) {
	return nil, nil
}

func (t track) Handler(message *discordgo.MessageCreate, matches []string, sendMessage func(message string)) *ThingInfo {
	id := t.Pattern().SubexpIndex("id")
	trackId := matches[id]

	res, err := t.spotify.client.GetTrack(t.spotify.ctx, spotify.ID(trackId))
	if err != nil {
		log.Println(fmt.Errorf("wow, got an error getting the Spotify track: %w", err))
	}
	sendMessage(fmt.Sprintf("This is a %s!", t.Name()))
	sendMessage("Found matching song!")
	sendMessage(fmt.Sprintf("This is %s!", res.Name))

	return &ThingInfo{
		Artist: res.Artists[0].Name,
		Name:   res.Name,
		Type:   t.HandlerType(),
		Link:   message.Content,
	}
}

func (track) HandlerType() ThingType {
	return ThingType("track")
}

func (track) Name() string {
	return "Spotify (track)"
}

func (track) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`https://open\.spotify\.com/track/(?P<id>[a-zA-Z0-9]+)`)
}

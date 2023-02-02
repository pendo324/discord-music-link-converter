package main

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

type album struct {
	spotify *spotifyClient
}

var _ Player = (*album)(nil)

func NewSpotifyAlbum(client *spotifyClient) *album {
	return &album{
		spotify: client,
	}
}

func (a album) Search(name string, artist string, thingType ThingType) *ThingInfo {
	return nil
}

func (a album) Handler(message *discordgo.MessageCreate, matches []string, sendMessage func(message string)) *ThingInfo {
	sendMessage(fmt.Sprintf("This is a %s!", a.Name()))
	return nil
}

func (album) HandlerType() ThingType {
	return ThingType("album")
}

func (album) Name() string {
	return "Spotify (album)"
}

func (album) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`https://open\.spotify\.com/album/(?P<id>[a-zA-Z0-9]+)`)
}

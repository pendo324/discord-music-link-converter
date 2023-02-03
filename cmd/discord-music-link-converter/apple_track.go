package main

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

type apple_track struct {
	apple *appleClient
}

var _ Player = (*apple_track)(nil)

func NewAppleTrack(client *appleClient) *apple_track {
	return &apple_track{
		apple: client,
	}
}

func (a apple_track) Search(name string, artist string, thingType ThingType) (*ThingInfo, error) {
	term := fmt.Sprintf("%s %s", name, artist)
	types := []string{"songs"}
	res, err := a.apple.Search(term, types)

	if err != nil {
		return nil, fmt.Errorf("failed to search apple music with term (%s) and types (%s): %w", term, types, err)
	}

	if len(res.Results.Songs.Data) > 0 {
		return &ThingInfo{
			Link:   res.Results.Songs.Data[0].Attributes.URL,
			Type:   ThingType("album"),
			Artist: res.Results.Songs.Data[0].Attributes.ArtistName,
			Name:   res.Results.Songs.Data[0].Attributes.Name,
		}, nil
	}

	return nil, fmt.Errorf("no results")
}

func (a apple_track) Handler(message *discordgo.MessageCreate, matches []string, sendMessage func(message string)) *ThingInfo {
	sendMessage(fmt.Sprintf("This is a %s!", a.Name()))
	return nil
}

func (apple_track) HandlerType() ThingType {
	return ThingType("track")
}

func (apple_track) Name() string {
	return "Apple (track)"
}

func (apple_track) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`https://music\.apple\.com/(?P<storefront>[a-z]+)/song/[0-9a-z-]+/(?P<id>[0-9]+)`)
}

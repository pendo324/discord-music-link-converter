package main

import (
	"fmt"
	"log"
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
			Type:   a.HandlerType(),
			Artist: res.Results.Songs.Data[0].Attributes.ArtistName,
			Name:   res.Results.Songs.Data[0].Attributes.Name,
		}, nil
	}

	return nil, fmt.Errorf("no results")
}

func (a apple_track) Handler(message *discordgo.MessageCreate, matches []string) *ThingInfo {
	id := a.Pattern().SubexpIndex("albumSongId")
	if id == -1 {
		// check if songId matches for alternate pattern
		id = a.Pattern().SubexpIndex("songId")
	}
	songId := matches[id]

	res, err := a.apple.GetSongById(songId)
	if err != nil {
		log.Println(fmt.Errorf("error getting the Apple track: %w", err))
	}

	return &ThingInfo{
		Artist: res.Attributes.ArtistName,
		Name:   res.Attributes.Name,
		Type:   a.HandlerType(),
		Link:   message.Content,
	}
}

func (apple_track) HandlerType() ThingType {
	return ThingType("track")
}

func (apple_track) Name() string {
	return "Apple (track)"
}

func (apple_track) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`https://music\.apple\.com/(?P<storefront>[a-z]+)/(album/[0-9a-z-]+/(?P<albumId>[0-9]+)\?i=(?P<albumSongId>[0-9]+)|(song/[0-9a-z-]+/(?P<songId>[0-9]+)))`)
}

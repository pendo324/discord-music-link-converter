package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

type apple_album struct {
	apple *appleClient
}

var _ Player = (*apple_album)(nil)

func NewAppleAlbum(client *appleClient) *apple_album {
	return &apple_album{
		apple: client,
	}
}

func (a apple_album) Search(name string, artist string, thingType ThingType) (*ThingInfo, error) {
	term := fmt.Sprintf("%s %s", name, artist)
	types := []string{"albums"}
	res, err := a.apple.Search(term, types)

	if err != nil {
		return nil, fmt.Errorf("failed to search apple music with term (%s) and types (%s): %w", term, types, err)
	}

	if len(res.Results.Albums.Data) > 0 {
		return &ThingInfo{
			Link:   res.Results.Albums.Data[0].Href,
			Type:   ThingType("album"),
			Artist: res.Results.Albums.Data[0].Attributes.ArtistName,
			Name:   res.Results.Albums.Data[0].Attributes.Name,
		}, nil
	}

	return nil, fmt.Errorf("no results")
}

func (a apple_album) Handler(message *discordgo.MessageCreate, matches []string, sendMessage func(message string)) *ThingInfo {
	id := a.Pattern().SubexpIndex("id")
	albumId := matches[id]

	res, err := a.apple.GetAlbumById(albumId)
	if err != nil {
		log.Println(fmt.Errorf("wow, got an error getting the Apple album: %w", err))
	}

	sendMessage(fmt.Sprintf("This is a %s!", a.Name()))
	sendMessage("Found matching album!")
	sendMessage(fmt.Sprintf("This is %s!", res.Attributes.Name))

	return &ThingInfo{
		Artist: res.Attributes.ArtistName,
		Name:   res.Attributes.Name,
		Type:   a.HandlerType(),
		Link:   message.Content,
	}
}

func (apple_album) HandlerType() ThingType {
	return ThingType("album")
}

func (apple_album) Name() string {
	return "Apple (album)"
}

func (apple_album) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`https://music\.apple\.com/(?P<storefront>[a-z]+)/album/[0-9a-z-]+/(?P<id>[0-9]+)`)
}

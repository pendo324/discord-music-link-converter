package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

type ThingType string

const (
	ALBUM ThingType = "album"
	TRACK ThingType = "track"
)

type ThingInfo struct {
	Artist string
	Name   string
	Link   string
	Type   ThingType
}

type Player interface {
	Search(name string, artist string, thingType ThingType) *ThingInfo
	Handler(message *discordgo.MessageCreate, matches []string, sendMessage func(message string)) *ThingInfo
	Name() string
	Pattern() *regexp.Regexp
	HandlerType() ThingType
}

func main() {
	botPtr := flag.String("disc-token", "", "The Discord bot token")
	spotCid := flag.String("spot-cid", "", "The Spotify client ID")
	spotSec := flag.String("spot-sec", "", "The Spotify client ID")

	flag.Parse()

	disc, err := NewDiscord(botPtr)
	if err != nil {
		log.Fatalf("Cannot setup discord client: %v", err)
	}

	spotifyClient, err := NewSpotifyClient(*spotCid, *spotSec)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to setup Spotify client: %w", err))
	}

	disc.ListenToMessages([]Player{NewSpotifyAlbum(spotifyClient), NewSpotifyTrack(spotifyClient)})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
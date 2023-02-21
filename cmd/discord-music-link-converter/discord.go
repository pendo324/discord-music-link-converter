package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Discord interface {
	ListenToMessages(p []Player)
}

var _ Discord = (*disc)(nil)

type disc struct {
	botToken *string
	session  *discordgo.Session
	state    *bool
}

func NewDiscord(botToken *string) (Discord, error) {
	// Tried to make fancy status thing, but its probably useless
	var mu sync.Mutex
	var state *bool
	sesh, _ := discordgo.New(fmt.Sprintf("Bot %s", *botToken))
	sesh.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		mu.Lock()
		fmt.Println("Bot is ready")
		isConnected := true
		state = &isConnected
		mu.Unlock()
	})

	sesh.AddHandler(func(s *discordgo.Session, r *discordgo.Disconnect) {
		mu.Lock()
		fmt.Println("Bot disconnected")
		isConnected := false
		state = &isConnected
		mu.Unlock()
	})

	err := sesh.Open()
	if err != nil {
		return nil, fmt.Errorf("err")
	}
	fmt.Println("Session opened")
	defer sesh.Close()

	return &disc{
		botToken: botToken,
		session:  sesh,
		state:    state,
	}, nil
}

func (d *disc) ListenToMessages(players []Player) {
	fmt.Println("Setting up handler...")
	d.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		var thread *discordgo.Channel
		for idx, player := range players {
			matches := player.Pattern().FindStringSubmatch(m.Content)
			if len(matches) > 0 {
				fmt.Printf("Original message in channel (%s) on server (%s): %s\n", m.Content, m.ChannelID, m.GuildID)

				if ch, err := s.State.Channel(m.ChannelID); err != nil || !ch.IsThread() {
					// get all other players of type and search for thing
					var otherHandlersOfType []Player
					for innerInx, innerPlayer := range players {
						if innerInx != idx && innerPlayer.HandlerType() == player.HandlerType() {
							otherHandlersOfType = append(otherHandlersOfType, innerPlayer)
						}
					}

					thingInfo := player.Handler(m, matches)

					if thread == nil {
						thread, err = s.MessageThreadStartComplex(m.ChannelID, m.ID, &discordgo.ThreadStart{
							Name:                fmt.Sprintf("Alternate links for %s: %s by %s", thingInfo.Type, thingInfo.Artist, thingInfo.Name),
							AutoArchiveDuration: 60,
							Invitable:           true,
						})
					}

					if err != nil {
						fmt.Println("failed to create thread: %w", err)
					}

					var embeds []*discordgo.MessageEmbed
					if thingInfo != nil {
						for _, p := range otherHandlersOfType {
							res, err := p.Search(thingInfo.Name, thingInfo.Artist, thingInfo.Type)
							if err != nil {
								log.Print("got error when searching: %w", err)
							}

							if res != nil {
								embeds = append(embeds, &discordgo.MessageEmbed{
									URL:   res.Link,
									Type:  "link",
									Title: p.Name(),
								})
							}
						}
					}

					// send message to thread with all the search info embedded
					_, err = s.ChannelMessageSendComplex(thread.ID, &discordgo.MessageSend{
						Embeds:  embeds,
						Content: "Found these links in other services",
					})
					if err != nil {
						log.Print("got error when sending reply in thread: %w", err)
					}
				}
			}
		}
	})
}

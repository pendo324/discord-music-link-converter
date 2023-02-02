package main

import (
	"bufio"
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/pendo324/discord-music-link-converter/pkg/util"
)

type appleClient struct{}

func NewAppleClient() (*appleClient, error) {
	req, err := http.NewRequest("GET", "https://music.apple.com/", strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var ctx context.Context
	var buf []byte

	allMatches := make(chan string)

	// TODO(Otto): idk help ty
	re := regexp.MustCompile(`[^"]*index.[a-z0-9]*.js`)

	go func() {
		defer close(allMatches)

		scanner := bufio.NewScanner(res.Body)
		scanner.Buffer(buf, 100)
		scanner.Split(util.SplitRegex(re))
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case allMatches <- scanner.Text():
			}
		}
	}()

	return nil, nil
}

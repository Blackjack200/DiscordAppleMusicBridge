package main

import (
	"fmt"
	"github.com/Blackjack200/DiscordAppleMusicBridge/applemusic"
	"github.com/hugolgst/rich-go/client"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := logrus.New()
	if err := client.Login("909326302757126186"); err != nil {
		log.Fatal(fmt.Errorf("failed to connect to discord client: %v", err))
	}
	log.Info("Discord <=> Apple Music Started!!!")
	ticker := time.NewTicker(time.Second * 2)
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-sigs
		ticker.Stop()
		client.Logout()
		log.Info("Discord <=> Apple Music Shutdown!!!")
		os.Exit(0)
	}()
	for {
		select {
		case <-ticker.C:
			if !update() {
				log.Error("failed to update")
			}
		}
	}
}

func update() bool {
	fetch, err := applemusic.Fetch()
	if err != nil {
		return client.SetActivity(client.Activity{
			Details:    "Idle",
			LargeImage: "apple_music_icon",
			LargeText:  "Apple Music",
		}) == nil
	}
	return client.SetActivity(client.Activity{
		Details:    fetch.Name,
		State:      fetch.Album,
		LargeImage: "apple_music_icon",
		LargeText:  fetch.Artist,
		Buttons: []*client.Button{
			{
				Label: fmt.Sprintf("Quailty: %vkhz/%vkbps", fetch.SampleRate/1000, fetch.BitRate),
				Url:   "https://music.apple.com",
			}, {
				Label: fmt.Sprintf("Disc: %v/%v Track: %v/%v", fetch.DiscNumber, fetch.DiscCount, fetch.TrackNumber, fetch.TrackCount),
				Url:   "https://music.apple.com",
			},
		},
	}) == nil
}

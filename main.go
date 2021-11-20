package main

import (
	"fmt"
	"github.com/Blackjack200/DiscordAppleMusicBridge/applemusic"
	"github.com/blackjack200/rich-go-plus/client"
	"github.com/blackjack200/rich-go-plus/codec"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var conn *client.Client

func reconnect() bool {
	var err error
	i := 16
	for i > 0 {
		i--
		c, e := client.Dial("909326302757126186")
		if e != nil {
			err = e
		} else {
			conn = c
			return true
		}
		time.Sleep(time.Second)
		logrus.Error(fmt.Errorf("failed to connect to discord client(retry %v): %v", i, err))
	}
	if err != nil {
		logrus.Fatal(fmt.Errorf("failed to connect to discord client: %v", err))
	}
	return false
}

func main() {
	log := logrus.New()
	reconnect()
	log.Info("Discord <=> Apple Music Started!!!")
	ticker := time.NewTicker(time.Second * 2)
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-sigs
		ticker.Stop()
		conn.Close()
		log.Info("Discord <=> Apple Music Shutdown!!!")
		os.Exit(0)
	}()
	for {
		select {
		case <-ticker.C:
			if err := update(); err != nil {
				log.Errorf("failed to update: %v", err)
				if reconnect() {
					log.Info("Recovered")
				}
			}
		}
	}
}

func update() error {
	fetch, err := applemusic.Fetch()
	if err != nil {
		return conn.SetActivity(&codec.Activity{
			Details:    "Idle",
			LargeImage: "apple_music_icon",
			LargeText:  "Apple Music",
		})
	}
	return conn.SetActivity(&codec.Activity{
		Details:    fmt.Sprintf("[%v] %v", strings.ToUpper(fetch.PlayerState), fetch.Name),
		State:      fetch.Album,
		LargeImage: "apple_music_icon",
		LargeText:  fetch.Artist,
		Buttons: []*codec.Button{
			{
				Label: fmt.Sprintf("Quailty: %vkhz/%vkbps", float32(fetch.SampleRate)/1000, fetch.BitRate),
				Url:   "https://music.apple.com",
			}, {
				Label: fmt.Sprintf("Disc: %v/%v Track: %v/%v", fetch.DiscNumber, fetch.DiscCount, fetch.TrackNumber, fetch.TrackCount),
				Url:   "https://music.apple.com",
			},
		},
	})
}

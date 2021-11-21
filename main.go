package main

import (
	"fmt"
	"github.com/Blackjack200/DiscordAppleMusicBridge/applemusic"
	"github.com/blackjack200/rich-go-plus/client"
	"github.com/blackjack200/rich-go-plus/codec"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/objc"
	"os"
	"runtime"
	"strings"
	"time"
)

var conn *client.Client

func reconnect() bool {
	i := 16
	for i > 0 {
		i--
		c, e := client.Dial("909326302757126186")
		if e == nil {
			conn = c
			return true
		}
		time.Sleep(time.Second)
	}
	return false
}

func main() {
	runtime.LockOSThread()
	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		setup()
	})
	app.Run()
}

func setup() {
	obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
	obj.Retain()
	obj.Button().SetTitle("🎵 Discord")

	ticker := time.NewTicker(time.Second * 2)
	go func() {
		reconnect()
		for {
			select {
			case <-ticker.C:
				if err := update(); err != nil {
				}
			}
		}
	}()

	itemQuit := cocoa.NSMenuItem_New()
	itemQuit.SetTitle("Quit")
	itemQuit.SetAction(objc.Sel("quit:"))
	cocoa.DefaultDelegateClass.AddMethod("quit:", func(_ objc.Object) {
		if conn != nil {
			conn.Close()
		}
		ticker.Stop()
		os.Exit(0)
	})

	menu := cocoa.NSMenu_New()
	menu.AddItem(itemQuit)
	obj.SetMenu(menu)
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

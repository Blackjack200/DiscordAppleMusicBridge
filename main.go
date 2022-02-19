package main

import (
	"github.com/Blackjack200/DiscordAppleMusicBridge/producer/applemusic"
	"github.com/progrium/macdriver/core"
	"os"
	"runtime"
	"time"

	"github.com/Blackjack200/DiscordAppleMusicBridge/producer"
	"github.com/Blackjack200/DiscordAppleMusicBridge/producer/vlc"
	"github.com/blackjack200/rich-go-plus/client"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/objc"
)

var conn *client.Client
var prod producer.Producer
var changed bool

func reconnect() bool {
	i := 16
	for i > 0 {
		i--
		c, _ := client.Dial(prod.AppId())
		if c != nil {
			conn = c
			return true
		}
		time.Sleep(time.Second / 4)
	}
	return false
}

func main() {
	runtime.LockOSThread()

	cocoa.TerminateAfterWindowsClose = false
	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		setup()
	})
	app.Run()
}

func setup() {
	obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
	obj.Retain()
	obj.Button().SetTitle("Discord")

	ticker := time.NewTicker(time.Second * 2)
	prod = vlc.NewProducer()
	go func() {
		reconnect()
		for {
			if changed {
				if conn != nil {
					conn.Close()
				}
				reconnect()
				changed = false
			}
			select {
			case <-ticker.C:
				if err := conn.SetActivity(prod.Activity()); err != nil {
					reconnect()
				}
			}
		}
	}()

	itemAppleMusic := cocoa.NSMenuItem_New()
	itemVLC := cocoa.NSMenuItem_New()
	itemQuit := cocoa.NSMenuItem_New()
	itemQuit.SetTitle("Quit")
	cocoa.DefaultDelegateClass.AddMethod("quit:", func(_ objc.Object) {
		if conn != nil {
			conn.Close()
		}
		ticker.Stop()
		os.Exit(0)
	})
	itemQuit.SetAction(objc.Sel("quit:"))

	itemAppleMusic.SetTitle("Apple Music")
	selectAppleMusicFunc := func(_ objc.Object) {
		selectAppleMusic()
		core.Dispatch(func() {
			itemVLC.SetEnabled(true)
			itemAppleMusic.SetEnabled(false)
			itemVLC.SetTitle("VLC")
			itemAppleMusic.SetTitle("Apple Music (selected)")
		})
	}
	cocoa.DefaultDelegateClass.AddMethod("a:", selectAppleMusicFunc)
	itemAppleMusic.SetAction(objc.Sel("a:"))

	itemVLC.SetTitle("VLC")
	cocoa.DefaultDelegateClass.AddMethod("vlc:", func(_ objc.Object) {
		selectVLC()
		core.Dispatch(func() {
			itemVLC.SetEnabled(false)
			itemAppleMusic.SetEnabled(true)
			itemVLC.SetTitle("VLC (selected)")
			itemAppleMusic.SetTitle("Apple Music")
		})
	})
	itemVLC.SetAction(objc.Sel("vlc:"))

	menu := cocoa.NSMenu_New()
	menu.AddItem(itemAppleMusic)
	menu.AddItem(itemVLC)
	menu.AddItem(itemQuit)

	obj.SetMenu(menu)
	selectAppleMusicFunc(nil)
}

func change(p producer.Producer) {
	if p != prod {
		prod = p
		changed = true
	}
}

func selectAppleMusic() {
	change(applemusic.NewProducer())
}

func selectVLC() {
	change(vlc.NewProducer())
}

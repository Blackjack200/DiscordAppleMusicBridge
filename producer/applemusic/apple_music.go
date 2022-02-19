package applemusic

import (
	_ "embed"
	"fmt"
	"strings"

	producer2 "github.com/Blackjack200/DiscordAppleMusicBridge/producer"
	"github.com/blackjack200/rich-go-plus/codec"
)

type producer struct{}

func (p producer) AppId() string {
	return "909326302757126186"
}

func (producer) Activity() *codec.Activity {
	d, err := fetch()
	if err != nil {
		return &codec.Activity{
			Details:    "Idle",
			LargeImage: "apple_music_icon",
			LargeText:  "Apple Music",
		}
	}
	return &codec.Activity{
		Details:    fmt.Sprintf("[%v] %v", strings.ToUpper(d.PlayerState), d.Name),
		State:      d.Album,
		LargeImage: "apple_music_icon",
		LargeText:  d.Artist,
		Buttons: []*codec.Button{
			{
				Label: fmt.Sprintf("Quailty: %vkhz/%vkbps", float32(d.SampleRate)/1000, d.BitRate),
				Url:   "https://music.apple.com",
			}, {
				Label: fmt.Sprintf("Disc: %v/%v Track: %v/%v", d.DiscNumber, d.DiscCount, d.TrackNumber, d.TrackCount),
				Url:   "https://music.apple.com",
			},
		},
	}
}

func NewProducer() producer2.Producer {
	return producer{}
}

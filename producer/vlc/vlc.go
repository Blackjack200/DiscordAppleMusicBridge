package vlc

import (
	"fmt"
	producer2 "github.com/Blackjack200/DiscordAppleMusicBridge/producer"
	"github.com/blackjack200/rich-go-plus/codec"
)

type producer struct{}

func (p producer) AppId() string {
	return "944543427658391572"
}

func (producer) Activity() *codec.Activity {
	t, err := fetch()
	if err != nil {
		return &codec.Activity{
			Details:    "Error",
			LargeImage: "vlc_icon",
			LargeText:  "VLC",
		}
	}
	isPlayingVideo := func() bool {
		return len(t.NameOfCurrentItem) != 0
	}
	if !isPlayingVideo() {
		return &codec.Activity{
			Details:    "Idle",
			LargeImage: "vlc_icon",
			LargeText:  "VLC",
		}
	}
	return &codec.Activity{
		Details:    fmt.Sprintf("[VLC] Playing"),
		State:      t.NameOfCurrentItem,
		LargeImage: "vlc_icon",
		LargeText:  "VLC",
		Buttons: []*codec.Button{
			{
				Label: fmt.Sprintf("Position: %v/%vs", t.CurrentTime, t.DurationOfCurrentItem),
				Url:   "https://www.videolan.org/",
			}, {
				Label: fmt.Sprintf("FullScreenMode: %v", t.FullscreenMode),
				Url:   "https://www.videolan.org/",
			},
		},
	}
}

func NewProducer() producer2.Producer {
	return producer{}
}

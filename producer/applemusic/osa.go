package applemusic

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type track struct {
	Name           string  `json:"name"`
	Kind           string  `json:"kind"`
	Album          string  `json:"album"`
	Artist         string  `json:"artist"`
	BitRate        int     `json:"bitRate"`
	DiscCount      int     `json:"discCount"`
	DiscNumber     int     `json:"discNumber"`
	Duration       float64 `json:"duration"`
	SampleRate     int     `json:"sampleRate"`
	TrackCount     int     `json:"trackCount"`
	TrackNumber    int     `json:"trackNumber"`
	PlayerState    string  `json:"playerState"`
	PlayerPosition float64 `json:"playerPosition"`
}

//go:embed fetch.js
var data []byte
var path string

func init() {
	path = filepath.Join(os.TempDir(), "discord_apple_music_bridge_fetch.js")
	if err := os.WriteFile(path, data, 0777); err != nil {
		logrus.Fatal(err)
	}
}

func fetch() (*track, error) {
	cmd := exec.Command("osascript", path)
	cmd.Env = os.Environ()

	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to read applemusic data: %v", err)
	} else {
		t := &track{}
		if err := json.Unmarshal(output, t); err != nil {
			return nil, fmt.Errorf("failed to marshal json data: %v", err)
		}
		return t, nil
	}
}

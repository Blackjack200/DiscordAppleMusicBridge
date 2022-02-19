package vlc

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
	Playing               bool   `json:"playing"`
	NameOfCurrentItem     string `json:"nameOfCurrentItem"`
	FullscreenMode        bool   `json:"fullscreenMode"`
	DurationOfCurrentItem int    `json:"durationOfCurrentItem"`
	CurrentTime           int    `json:"currentTime"`
}

//go:embed fetch.js
var data []byte
var path string

func init() {
	path = filepath.Join(os.TempDir(), "discord_vlc_bridge_fetch.js")
	if err := os.WriteFile(path, data, 0777); err != nil {
		logrus.Fatal(err)
	}
}

func fetch() (*track, error) {
	cmd := exec.Command("osascript", path)
	cmd.Env = os.Environ()

	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to read vlc data: %v", err)
	} else {
		t := &track{}
		if err := json.Unmarshal(output, t); err != nil {
			return nil, fmt.Errorf("failed to marshal json data: %v", err)
		}
		return t, nil
	}
}

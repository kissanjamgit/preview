// Package config to load toml
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Player     string `toml:"player"`
	PlayerArgs string `toml:"playerArgs"`
}

func New() (c Config, err error) {
	c = Config{Player: "vlc.exe", PlayerArgs: ""}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configPath := filepath.Join(homeDir, ".preview.toml")
	b, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return c, nil
	}

	if err != nil {
		return
	}
	var t Config

	if err = toml.Unmarshal(b, &t); err != nil {
		return
	} else if (Config{}) != t { // if config is empty
		c = t
	}
	return
}

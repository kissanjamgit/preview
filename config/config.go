// Package config to load toml
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ConfigPlayer struct {
	Player     string `toml:"player"`
	PlayerArgs string `toml:"playerArgs"`
}
type ConfigNetwork struct {
	Proxy string `toml:"proxy"`
}

type Config struct {
	ConfigPlayer
	ConfigNetwork
}

var ConfigLazy *Config

var (
	configPlayer  = ConfigPlayer{Player: "vlc.exe", PlayerArgs: ""}
	configNetwork = ConfigNetwork{Proxy: `http://127.0.0.1:18080`}
)

func New() (c Config, err error) {
	c = Config{ConfigPlayer: configPlayer, ConfigNetwork: configNetwork}
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

	if err = toml.Unmarshal(b, &c); err != nil {
		return
	}
	if (c.ConfigPlayer == ConfigPlayer{}) {
		c.ConfigPlayer = configPlayer
	}
	if (c.ConfigNetwork == ConfigNetwork{}) {
		c.ConfigNetwork = configNetwork
	}
	ConfigLazy = &c

	return
}

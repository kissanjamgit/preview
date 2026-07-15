// Package common do something
package common

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/kissanjamgit/preview"
	"github.com/kissanjamgit/preview/config"
)

// NewClient creates a pre-configured resty client
func NewClient() *resty.Client {
	client := resty.New()

	// Apply proxy from config if available
	if config.ConfigLazy != nil && config.ConfigLazy.Proxy != "" {
		client.SetProxy(config.ConfigLazy.Proxy)
	}

	// Set global default headers
	client.SetHeaders(map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:152.0) Gecko/20100101 Firefox/152.0",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9",
	})

	return client
}

func PR2String(pr []preview.ContentResource) (string, error) {
	if len(pr) == 0 {
		return ``, fmt.Errorf("len(pr) == 0")
	}
	var buffer strings.Builder
	buffer.WriteString("#EXTM3U\n")
	for _, cr := range pr {
		fmt.Fprintf(&buffer, "#EXTINF:-1,%s\n%s\n", cr.Source, cr.View)
	}
	return buffer.String(), nil
}

var PlayerArgs = []string{"-"}

func Play(cfg config.Config, str string) (err error) {
	PlayerArgs = append(PlayerArgs, strings.Split(cfg.PlayerArgs, " ")...)
	playCmd := exec.Command(cfg.Player, PlayerArgs...)
	stdin, err := playCmd.StdinPipe()
	if err != nil {
		return
	}
	err = playCmd.Start()
	if err != nil {
		return
	}
	_, err = stdin.Write([]byte(str))
	if err != nil {
		return
	}
	err = stdin.Close()
	if err != nil {
		return
	}
	err = playCmd.Wait()

	return
}

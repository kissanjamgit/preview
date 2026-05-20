// Package common do something
package common

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/kissanjamgit/preview"
	"github.com/kissanjamgit/preview/config"
)

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

func Play(cfg config.Config, str string) (err error) {
	PlayerArgs := []string{"-"} // stdin input
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

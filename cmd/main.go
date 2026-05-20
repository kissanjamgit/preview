package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kissanjamgit/preview/barrel"
	"github.com/kissanjamgit/preview/common"
	"github.com/kissanjamgit/preview/config"
	"github.com/kissanjamgit/preview/search"

	"github.com/kissanjamgit/preview"

	"github.com/spf13/cobra"
)

var example = func() string {
	var buff strings.Builder
	for i, p := range barrel.Domain {
		fmt.Fprintf(&buff, "%2d: %s\n", i, p.Name())
	}
	return buff.String()
}()

func getM3U(cmd *cobra.Command, args []string, index int) (m3u string, err error) {
	if len(args) == 0 {
		err = fmt.Errorf("need at least one argument")
		return
	}

	input, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	view := barrel.Domain[input]
	if view, ok := view.(preview.Domain); ok {
		var buff strings.Builder
		for i, d := range view.Domain() {
			fmt.Fprintf(&buff, "%2d: %s\n", i, d)
		}
		cmd.Example = buff.String()

		if len(args) < 2 {
			err = fmt.Errorf("need one more argument")
			return
		}
		subIndex, err := strconv.Atoi(args[1])
		if err != nil {
			return ``, err
		}
		if subIndex < 0 || subIndex >= len(view.Domain()) {
			err = fmt.Errorf("subIndex >= 0 && subIndex < len(f.List); subindex: %d", subIndex)
			return ``, err
		}
		view.SetSource(view.Domain()[subIndex])
	}

	pr, err := view.Get(index)
	if err != nil {
		return ``, err
	}
	m3u, err = common.PR2String(pr)
	if err != nil {
		return ``, err
	}
	return
}

func cli() (err error) {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cbr := cobra.Command{
		Use: "view",
	}
	var index int
	show := cobra.Command{
		Use:     "show",
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			m3u, err := getM3U(cmd, args, index)
			if err != nil {
				return err
			}
			fmt.Println(m3u)
			return err
		},
	}

	play := cobra.Command{
		Use:     "play",
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			m3u, err := getM3U(cmd, args, index)
			if err != nil {
				return err
			}

			return common.Play(cfg, m3u)
		},
	}

	show.Flags().IntVar(&index, "index", 0, "index value")
	play.Flags().IntVar(&index, "index", 0, "index value")
	search.Add(&cbr, cfg)
	cbr.AddCommand(&show)
	cbr.AddCommand(&play)
	err = cbr.Execute()
	if err != nil {
		os.Exit(1)
	}
	return
}

func main() {
	cli()
}

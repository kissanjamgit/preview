package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kissanjamgit/preview/barrel"
	"github.com/kissanjamgit/preview/config"

	"github.com/kissanjamgit/preview"

	"github.com/spf13/cobra"
)

func PR2String(pr []preview.ContentResource) (string, error) {
	if len(pr) == 0 {
		return ``, fmt.Errorf("len(pr) == 0")
	}
	var buffer strings.Builder
	buffer.WriteString("#EXTM3U\n")
	for _, cr := range pr {
		buffer.WriteString("#EXTINF:-1," + cr.Source + "\n" + cr.View + "\n")
	}
	return buffer.String(), nil
}

var exampleSearch = func() string {
	var buff strings.Builder
	listLength := 0
	for _, p := range barrel.Domain {
		search, ok := p.(preview.Search)
		if !ok {
			continue
		}

		fmt.Fprintf(&buff, "%2d: %s\n", listLength, search.Name())
		listLength += 1
	}
	return buff.String()
}()

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
	m3u, err = PR2String(pr)
	if err != nil {
		return ``, err
	}
	return
}

func search(cfg config.Config, index *int) *cobra.Command {
	return &cobra.Command{
		Use:     "search",
		Example: exampleSearch,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) == 0 {
				err = fmt.Errorf("need at least one argument")
				return
			}
			searchList := []preview.Search{}
			for _, item := range barrel.Domain {
				p, ok := item.(preview.Search)
				if !ok {
					continue
				}
				searchList = append(searchList, p)
			}
			listIndex, err := strconv.Atoi(args[0])
			if err != nil {
				return
			}
			if len(searchList) < listIndex {
				return fmt.Errorf(`len(searchList) < listIndex`)
			}

			var pr []preview.ContentResource
			search := searchList[listIndex]
			if searchDomain, ok := search.(preview.SearchDomain); ok {
				var buff strings.Builder
				for i, d := range searchDomain.Domain() {
					fmt.Fprintf(&buff, "%2d: %s\n", i, d)
				}

				cmd.Example = buff.String()

				if len(args) < 3 {
					err = fmt.Errorf("len(args) < 3")
					return
				}
				sublistIndex, err := strconv.Atoi(args[1])
				if err != nil {
					return err
				}
				searchDomain.SetSource(searchDomain.Domain()[sublistIndex])
				pr, err = searchDomain.Search(args[2:], *index)
				if err != nil {
					return err
				}
			} else {
				pr, err = search.Search(args[1:], *index)
				if err != nil {
					return err
				}
			}
			// var hostname string
			// if uri, err := url.Parse(args[0]); err != nil && uri.Hostname() != `` {
			// 	return err
			// } else {
			// 	hostname = uri.Hostname()
			// }
			//
			// switch {
			// case strings.Contains(hostname, `manyvids`):
			// 	m := manyv.New(args[0])
			// 	m.Search(index)
			// default:
			// 	err = fmt.Errorf("switch case exhausted")
			// 	return
			//
			// }

			m3u, err := PR2String(pr)
			if err != nil {
				return err
			}
			switch cmd.Parent().Name() {
			case `show`:
				fmt.Println(m3u)
			case `play`:
				err = cobraPlay(cfg, m3u)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("switch case exhausted")
			}
			return
		},
	}
}

func cobraPlay(cfg config.Config, str string) (err error) {
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

			return cobraPlay(cfg, m3u)
		},
	}

	show.Flags().IntVar(&index, "index", 0, "index value")
	play.Flags().IntVar(&index, "index", 0, "index value")
	cbr.AddCommand(&show)
	cbr.AddCommand(&play)
	searchPlay := search(cfg, &index)
	searchPlay.Flags().IntVar(&index, "index", 0, "index value")
	searchShow := search(cfg, &index)
	searchShow.Flags().IntVar(&index, "index", 0, "index value")
	play.AddCommand(searchPlay)
	show.AddCommand(searchShow)
	err = cbr.Execute()
	if err != nil {
		os.Exit(1)
	}
	return
}

func main() {
	cli()
}

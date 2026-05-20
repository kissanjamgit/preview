// Package search provides preview
package search

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kissanjamgit/preview"
	"github.com/kissanjamgit/preview/barrel"
	"github.com/kissanjamgit/preview/common"
	"github.com/kissanjamgit/preview/config"
	"github.com/spf13/cobra"
)

var search = cobra.Command{
	Use:     "search",
	Example: exampleSearch,
}

var (
	index         int
	cfg           config.Config
	exampleSearch = func() string {
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
)

func Add(cbr *cobra.Command, c config.Config) {
	cfg = c
	searchPlay.Flags().IntVar(&index, `index`, 0, `index`)
	searchShow.Flags().IntVar(&index, `index`, 0, `index`)
	search.AddCommand(&searchPlay)
	search.AddCommand(&searchShow)
	cbr.AddCommand(&search)
}

var searchShow = cobra.Command{
	Use:     "show",
	Example: exampleSearch,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		m3u, err := getM3UFromSearch(cmd, args, index)
		if err != nil {
			return
		}
		fmt.Println(m3u)
		return
	},
}

var searchPlay = cobra.Command{
	Use:     "play",
	Example: exampleSearch,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		m3u, err := getM3UFromSearch(cmd, args, index)
		if err != nil {
			return
		}
		common.Play(cfg, m3u)
		return
	},
}

func getM3UFromSearch(cmd *cobra.Command, args []string, index int) (m3u string, err error) {
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
		err = fmt.Errorf(`len(searchList) < listIndex`)
		return
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
			return ``, err
		}
		searchDomain.SetSource(searchDomain.Domain()[sublistIndex])
		pr, err = searchDomain.Search(args[2:], index)
		if err != nil {
			return ``, err
		}
	} else {

		if len(args) < 2 {
			err = fmt.Errorf("len(args) < 2")
			return
		}
		pr, err = search.Search(args[1:], index)
		if err != nil {
			return ``, err
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

	m3u, err = common.PR2String(pr)
	if err != nil {
		return
	}
	return
}

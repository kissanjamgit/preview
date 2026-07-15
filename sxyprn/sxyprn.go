package sxyprn

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kissanjamgit/preview"
	"github.com/kissanjamgit/preview/common"
)

type sxyprn struct {
	Source string
}

func New() *sxyprn {
	return &sxyprn{}
}

func (s *sxyprn) Name() string {
	return "sxyprn"
}

func (s *sxyprn) SetSource(source string) {
	s.Source = source
}

func (s *sxyprn) Get(index int) (list []preview.ContentResource, err error) {
	// Construct URL
	page := 30 * index
	uri := fmt.Sprintf("https://sxyprn.com/http.html?page=%d", page)

	client := common.NewClient()
	res, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36").
		Get(uri)

	if err != nil {
		return
	}

	if res.StatusCode() != 200 {
		err = fmt.Errorf("status code: %d", res.StatusCode())
		return
	}

	d, err := goquery.NewDocumentFromReader(strings.NewReader(res.String()))
	if err != nil {
		return
	}

	d.Find(".post_el_small").Each(func(i int, g *goquery.Selection) {
		aTag := g.Find("a.tdn.post_time")
		title, ok := aTag.Attr("title")
		if !ok {
			return
		}
		src, ok := g.Find(".hvp_player").Attr("src")
		if !ok {
			return
		}

		list = append(list, preview.ContentResource{
			Source: strings.NewReplacer("\n", "", "\r", "").Replace(title),
			View:   "https:" + src,
		})
	})

	return
}

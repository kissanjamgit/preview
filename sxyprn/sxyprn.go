package sxyprn

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kissanjamgit/preview"
	"github.com/kissanjamgit/preview/config"
	"resty.dev/v3"
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

func (s *sxyprn) SearchIdentity() {}

func (s *sxyprn) Search(query []string, index int) ([]preview.ContentResource, error) {
	page := 30 * index
	// Construct search URL based on query
	q := strings.Join(query, "-")
	uri := fmt.Sprintf("https://sxyprn.com/%s.html?page=%d", url.PathEscape(q), page)
	return s.fetch(uri)
}

func (s *sxyprn) Get(index int) (list []preview.ContentResource, err error) {
	page := 30 * index
	uri := fmt.Sprintf("https://sxyprn.com/http.html?page=%d", page)
	return s.fetch(uri)
}

func (s *sxyprn) fetch(uri string) (list []preview.ContentResource, err error) {
	client := resty.New()
	if config.ConfigLazy != nil && config.ConfigLazy.Proxy != "" {
		client.SetProxy(config.ConfigLazy.Proxy)
	}

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

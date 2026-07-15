// Package sxyprn provides a provider for https://sxyprn.com
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

// --- Shared Logic ---

func fetch(uri string) (*goquery.Document, error) {
	client := resty.New()
	if config.ConfigLazy != nil && config.ConfigLazy.Proxy != "" {
		client.SetProxy(config.ConfigLazy.Proxy)
	}

	res, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/139.0.0.0").
		Get(uri)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode())
	}

	return goquery.NewDocumentFromReader(strings.NewReader(res.String()))
}

// --- Provider 1: Main Site ---

type Sxyprn struct{}

func New() *Sxyprn                        { return &Sxyprn{} }
func (s *Sxyprn) Name() string            { return "sxyprn" }
func (s *Sxyprn) SetSource(source string) {}
func (s *Sxyprn) SearchIdentity()         {}

func (s *Sxyprn) Search(query []string, index int) ([]preview.ContentResource, error) {
	page := 30 * index
	q := strings.Join(query, "-")
	uri := fmt.Sprintf("https://sxyprn.com/%s.html?page=%d", url.PathEscape(q), page)
	return s.parse(uri)
}

func (s *Sxyprn) Get(index int) (list []preview.ContentResource, err error) {
	uri := fmt.Sprintf("https://sxyprn.com/http.html?page=%d", 30*index)
	return s.parse(uri)
}

func (s *Sxyprn) parse(uri string) (list []preview.ContentResource, err error) {
	d, err := fetch(uri)
	if err != nil {
		return
	}
	d.Find(".post_el_small").Each(func(i int, g *goquery.Selection) {
		aTag := g.Find("a.tdn.post_time")
		title, ok := aTag.Attr("title")
		src, ok2 := g.Find(".hvp_player").Attr("src")
		if ok && ok2 {
			list = append(list, preview.ContentResource{
				Source: strings.NewReplacer("\n", "", "\r", "").Replace(title),
				View:   "https:" + src,
			})
		}
	})
	return
}

// --- Provider 2: Blog Site ---

type SxyprnBlog struct{}

func NewBlog() *SxyprnBlog                    { return &SxyprnBlog{} }
func (s *SxyprnBlog) Name() string            { return "sxyprn-blog" }
func (s *SxyprnBlog) SetSource(source string) {}
func (s *SxyprnBlog) Get(index int) (list []preview.ContentResource, err error) {
	uri := fmt.Sprintf("https://sxyprn.net/blog/all/%d", 20*index)
	d, err := fetch(uri)
	if err != nil {
		return
	}

	d.Find(".post_el_small").Each(func(i int, g *goquery.Selection) {
		title := strings.TrimSpace(g.Find(".post_text").Text())
		// extLinkTag := g.Find(`.extlink_icon.extlink`)

		// if extLinkTag.Length() > 0 {
		// 	srcJpg, ok0 := g.Find(".mini_post_vid_thumb.lazyloaded").Attr("data-src")
		// 	href, ok1 := extLinkTag.Attr(`href`)
		// 	if ok0 && ok1 {
		// 		view := srcJpg
		// 		if strings.HasPrefix(srcJpg, "//") {
		// 			view = "https:" + srcJpg
		// 		}
		// 		list = append(list, preview.ContentResource{
		// 			Source: title + ` ` + href,
		// 			View:   view,
		// 		})
		// 		return
		// 	}
		// }

		src, ok := g.Find(".hvp_player").Attr("src")
		if ok {
			list = append(list, preview.ContentResource{
				Source: title,
				View:   "https:" + src,
			})
		}
	})
	return
}

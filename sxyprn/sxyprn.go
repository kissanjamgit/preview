// Package sxyprn provides a provider for https://sxyprn.com
package sxyprn

import (
	"fmt"
	"net/url"
	"slices"
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

var blog = map[string]string{`sxyprnBlog`: `all`, `ILUVY`: `5f3950a938042`, `AJ47`: `@AJ47`, `mikess`: `@mikess`, `Ajx`: `@Ajx`, `PornoBB`: `66326b4fdfd13`, `DemonINC`: `64d19fdf6970d`, `iAmFckr`: `667bfa40a223f`, `RIPNSFW`: `604868c9e12f6`, `sandy998`: `67b6a8c5c1227`, `scfc`: `6803728ea8837`, `VOOP`: `6a1f341252008`, `JORDAN888`: `643b5d34ef375`}

type SxyprnBlog struct {
	Source string
}

func NewBlog() *SxyprnBlog         { return &SxyprnBlog{} }
func (s *SxyprnBlog) Name() string { return "sxyprn-blog" }
func (s *SxyprnBlog) Domain() []string {
	keys := make([]string, len(blog))
	i := 0
	for k := range blog {
		keys[i] = k
		i++
	}
	slices.Sort(keys)
	return keys
}
func (s *SxyprnBlog) SetSource(source string) { s.Source = source }
func (s *SxyprnBlog) Get(index int) (list []preview.ContentResource, err error) {
	uri := fmt.Sprintf("https://sxyprn.net/blog/%s/%d", blog[s.Source], 20*index)
	d, err := fetch(uri)
	if err != nil {
		return
	}

	d.Find(".post_el_small").Each(func(i int, g *goquery.Selection) {
		// 1. Corrected selector with dots
		extlinkIcon := g.Find(`.extlink_icon.extlink`)

		var hrefs []string
		extlinkIcon.Each(func(i int, s *goquery.Selection) {
			text, ok := s.Attr(`href`)
			if ok {
				hrefs = append(hrefs, text)
			}
		})

		title := strings.TrimSpace(g.Find(".post_text").Text())
		title = strings.ReplaceAll(title, "\n", " ")

		// 2. Only add the links if we actually found any
		if len(hrefs) > 0 {
			title += " " + strings.Join(hrefs, " ")
		}
		var src string
		if g.Find(".duration_small").Text() == `EXTERNAL LINK` {
			if s, ok := g.Find(".mini_post_vid_thumb").Attr("data-src"); !ok {
				return
			} else {
				src = s
			}
		} else {
			if s, ok := g.Find(".hvp_player").Attr("src"); !ok {
				return
			} else {
				src = s
			}
		}
		list = append(list, preview.ContentResource{
			Source: title,
			View:   "https:" + src,
		})
	})
	return
}

// --- Provider 3: Popular Site ---

type SxyprnPopular struct{}

func NewPopular() *SxyprnPopular                 { return &SxyprnPopular{} }
func (s *SxyprnPopular) Name() string            { return "sxyprn-popular" }
func (s *SxyprnPopular) SetSource(source string) {}
func (s *SxyprnPopular) Get(index int) (list []preview.ContentResource, err error) {
	uri := fmt.Sprintf("https://sxyprn.net/popular/top-viewed/%d?p=day", 30*index)
	d, err := fetch(uri)
	if err != nil {
		return
	}

	d.Find(".post_el_small").Each(func(i int, g *goquery.Selection) {
		extlinkIcon := g.Find(`.extlink_icon.extlink`)

		var hrefs []string
		extlinkIcon.Each(func(i int, s *goquery.Selection) {
			text, ok := s.Attr(`href`)
			if ok {
				hrefs = append(hrefs, text)
			}
		})

		title := strings.TrimSpace(g.Find(".post_text").Text())
		title = strings.ReplaceAll(title, "\n", " ")

		if len(hrefs) > 0 {
			title += " " + strings.Join(hrefs, " ")
		}

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

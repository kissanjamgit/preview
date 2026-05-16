// Package hd8k do something
package hd8k

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"

	"resty.dev/v3"
)

type HD8k struct {
	source string
}

var domain = "https://en8.pornhd8k.me"

func (h HD8k) Name() string {
	return "hd8k"
}

func (h *HD8k) SetSource(source string) {
	h.source = source
}

func New() *HD8k {
	return &HD8k{}
}

func (h HD8k) Get(index int) (pr []preview.ContentResource, err error) {
	if index <= 0 {
		index = 1
	}
	client := resty.New()
	uri := fmt.Sprintf("https://en8.pornhd8k.me/porn-hd-videos/page-%d", index)
	// uri = "https://en8.pornhd8k.me/porn-hd-videos/page-1"
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	matchSource := regexp.MustCompile(`href="(/movies/[^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	matchView := regexp.MustCompile(`data-preview="([^"]*)"`).FindAllStringSubmatch(res.String(), -1)
	matchImage := regexp.MustCompile(`data-original="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)

	if len(matchSource) != len(matchView) || len(matchSource) != len(matchImage) {
		err = fmt.Errorf("length of source and view are not equale, len(matchSource): %d len(matchView): %d", len(matchSource), len(matchView))
		return
	}

	for i, source := range matchSource {
		view := matchView[i]
		if len(source) < 2 || len(view) < 2 {
			fmt.Println(source, view)
			continue
		}
		v := view[1]
		if v == `` || len(v) > 2 && v[:2] == "//" {
			if imageSubmatch := matchImage[i]; len(imageSubmatch) >= 2 {
				v = domain + matchImage[i][1]
			}
		}
		pr = append(pr, preview.ContentResource{View: v, Source: domain + source[1]})
	}

	return
}

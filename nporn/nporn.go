// Package nporn so something
package nporn

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type nporn struct {
	Source string
}

func New(source string) preview.Preview {
	return &nporn{source}
}

var size = 12

var baseURL = "https://nubiles-porn.com"

func (n *nporn) Get(index int) (list []preview.ContentResource, err error) {
	client := resty.New()
	uri := fmt.Sprintf("%s/video/gallery/%d", baseURL, size*index)
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`<a href="(/video/watch[^"]+)">\s*<div class="overlay-video-wrapper hover-thumb cover" data-preview-src="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, item := range submatch {
		if len(item) < 3 {
			continue
		}
		list = append(list, preview.ContentResource{Source: baseURL + item[1], View: html.UnescapeString(strings.TrimSpace(item[2]))})

	}
	return
}

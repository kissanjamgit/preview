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

func (n *nporn) SetSource(source string) {
	n.Source = source
}

func (n nporn) Name() string {
	return "nporn"
}

func New() *nporn {
	return &nporn{}
}

var size = 12

func (nporn) Domain() []string {
	return Domain
}

var Domain = []string{`nubiles-porn`, `brattysis`, `momlover`, `shesbreedingmaterial`, `realitysis`, `caughtmycoach`, `cheatingsis`, `cumswappingsis`, `myfamilypies`, `stepsiblingscaught`, `familyswap`, `momsteachsex`}

func (n nporn) Get(index int) (list []preview.ContentResource, err error) {
	client := resty.New()
	baseURL := `https://` + n.Source + `.com`
	uri := fmt.Sprintf("%s/video/gallery/%d", baseURL, size*index)
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`<a href="(/video/watch[^"]+)">\s*<div[^>]*data-preview-src="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, item := range submatch {
		if len(item) < 3 {
			continue
		}
		list = append(list, preview.ContentResource{Source: baseURL + item[1], View: html.UnescapeString(strings.TrimSpace(item[2]))})

	}
	return
}

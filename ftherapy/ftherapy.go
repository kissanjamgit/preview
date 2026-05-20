// Package ftherapy provides preview
package ftherapy

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type ftherapy struct {
	source string
}

func (f *ftherapy) SetSource(source string) {
	f.source = source
}

func (f ftherapy) Name() string {
	return "ftherapy"
}

func New() *ftherapy {
	return &ftherapy{}
}

func (f ftherapy) Domain() []string {
	return Domain
}

var Domain = []string{`familytherapyxxx`, `momcomesfirst`, `cockninjastudios`, `perfectgirlfriend`, `analtherapyxxx`, `wifelovesblack`, `teenlovesblack`}

type pathConv struct {
	name string
	img  string
	av   string
}

var tree = []pathConv{
	{`familytherapyxxx`, `p?\d*\.[a-z]+`, `web.mp4`},
	{`momcomesfirst`, `p?\d*\.[a-z]+`, `web.mp4`},
	{`cockninjastudios`, `thumb\.[a-z]+`, `preview.mp4`},
	{`perfectgirlfriend`, `p?\d*\.[a-z]+`, `web.mp4`},
	{`analtherapyxxx`, `p?\d*\.[a-z]+`, `web.mp4`},
	{`wifelovesblack`, `_thumb\.[a-z]+`, `_trailer.mp4`},
	{`teenlovesblack`, `-thumb\.[a-z]+`, `_trailer.mp4`},
}

func (f ftherapy) Get(index int) (pr []preview.ContentResource, err error) {
	sp, err := func() (pathConv, error) {
		for _, i := range tree {
			if i.name != f.source {
				continue
			}
			return i, nil
		}
		return pathConv{}, fmt.Errorf(`f.source not in tree`)
	}()
	index += 1
	uri := fmt.Sprintf(`https://%s.com/page/%d/?et_blog`, f.source, index)
	client := resty.New()
	defer client.Close()
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`<a href="([^"]+)" class="entry-featured-image-url"><img[^>]*src="([^"]+)" `).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatch {
		if len(m) < 2 {
			continue
		}

		split := strings.Split(strings.TrimRight(m[2], `/`), `/`)
		last := len(split) - 1
		split[last] = regexp.MustCompile(sp.img).ReplaceAllString(split[last], sp.av)
		view := strings.Join(split, `/`)
		pr = append(pr, preview.ContentResource{Source: m[1], View: view})
	}
	return
}

// Package newsens some
package newsens

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type newsens struct {
	Source string
}

func (n *newsens) SetSource(source string) {
	n.Source = source
}

func (n newsens) Name() string {
	return "newsenss"
}

func New() *newsens {
	return &newsens{}
}

type StudioPath struct {
	name string
	path string
}

var tree = []StudioPath{
	{`newsensations`, `tour_ns`},
	{`familyxxx`, `tour_famxxx`},
	{`hotwifexxx`, `tour_hwxxx`},
	{`girlgirlxxx`, `tour_girlgirlxxx`},
	{`shanedieselxxx`, `tour_sdxxx`},
}

var Domain = func() (s []string) {
	for _, i := range tree {
		s = append(s, i.name)
	}
	return
}()

func (n newsens) Get(index int) (list []preview.ContentResource, err error) {
	index += 1
	studio, err := func() (StudioPath, error) {
		for _, s := range tree {
			if n.Source != s.name {
				continue
			}
			return s, nil
		}
		return StudioPath{}, fmt.Errorf("n.Source not in tree")
	}()
	if err != nil {
		return
	}

	uri := fmt.Sprintf(`https://www.%s.com/%s/categories/movies_%d_d.html`, studio.name, studio.path, index) // 0 isn't allowed;
	client := resty.New()
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	listSrc := []string{}
	listHref := []string{}
	submatchSrc := regexp.MustCompile(`source src="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatchSrc {
		if len(m) < 2 {
			continue
		}
		listSrc = append(listSrc, m[1])
	}
	submatchHref := regexp.MustCompile(`<h4><a href="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatchHref {
		if len(m) < 2 {
			continue
		}
		listHref = append(listHref, m[1])
	}
	if len(listSrc) != len(listHref) {
		err = fmt.Errorf("len(listSrc) != len(listHref)")
		return
	}

	for index, src := range listSrc {
		list = append(list, preview.ContentResource{Source: listHref[index], View: src})
	}
	return
}

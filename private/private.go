// Package private some
package private

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type private struct {
	source string
}

func (p *private) Name() string {
	return "private"
}

func (p *private) SetSource(source string) {
	p.source = source
}

func New() *private {
	return &private{}
}

func (p private) Get(index int) (list []preview.ContentResource, err error) {
	index += 1

	uri := fmt.Sprintf(`https://www.private.com/scenes/%d/`, index) // 0 isn't allowed;
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
	submatchHref := regexp.MustCompile(`data-track="TITLE_LINK" href="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
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

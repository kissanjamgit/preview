// Package enjoyx some
package enjoyx

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type enjoyx struct {
	Source string
}

func (e enjoyx) Name() string {
	return "enjoyx"
}

func (e *enjoyx) SetSource(source string) {
	e.Source = source
}

func New() *enjoyx {
	return &enjoyx{}
}

func (e enjoyx) Get(index int) (list []preview.ContentResource, err error) {
	index += 1
	uri := fmt.Sprintf(`https://enjoyx.com/video?page=%d`, index)
	client := resty.New()
	defer client.Close()
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	listSrc := []string{}
	listHref := []string{}
	submatch := regexp.MustCompile(`<video[^>]+src="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatch {
		if len(m) < 2 {
			continue
		}
		listSrc = append(listSrc, m[1])
	}

	submatch = regexp.MustCompile(`class="video-card__item" href="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatch {
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

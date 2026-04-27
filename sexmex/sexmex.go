// Package smex is a package a
package smex

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"

	"resty.dev/v3"
)

type Smex struct {
	source string
}

func (s *Smex) Get(index int) (pr []preview.ContentResource, err error) {
	client := resty.New().R()
	res, err := client.Get("https://sexmex.xxx/tour/categories/movies.html")
	if err != nil {
		return
	}
	matchName := regexp.MustCompile(`https://sexmex.xxx/tour/updates/[^"]*`).FindAllString(res.String(), -1)

	submatchView := regexp.MustCompile(`clip directory ([^ -]*)`).FindAllStringSubmatch(res.String(), -1)

	matchView := []string{}

	for _, name := range submatchView {
		if len(name) < 2 {
			continue
		}
		matchView = append(matchView, fmt.Sprintf(`https://sexmex.xxx/trailers/%s.mp4`, name[1]))
	}
	if len(matchName) != 2*len(matchView) {
		err = fmt.Errorf("len(matchName) != 2*len(matchView)  len(matchName): %d len(matchView): %d", len(matchName), len(matchView))
		return
	}

	for i, view := range matchView {
		pr = append(pr, preview.ContentResource{Source: matchName[i*2], View: view})
	}
	return
}

func New(source string) preview.Preview {
	return &Smex{source}
}

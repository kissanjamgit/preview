// Package pworld something something
package pworld

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"

	"resty.dev/v3"
)

type PWorld struct {
	source string
}

func (p *PWorld) SetSource(source string) {
	p.source = source
}

func (p PWorld) Get(index int) (pr []preview.ContentResource, err error) {
	client := resty.New().R()
	url := fmt.Sprintf("https://pornworld.com/videos?page=%d", index)
	res, err := client.Execute(resty.MethodGet, url)
	if err != nil {
		return
	}
	matchSource := regexp.MustCompile(`pornworld.com/watch[^"]*`).FindAllString(string(res.Bytes()), -1)

	matchView := regexp.MustCompile(`https://[^"]*trailer.mp4`).FindAllString(string(res.Bytes()), -1)

	if len(matchSource) == 2*len(matchView) {
		l := []string{}
		for i := range len(matchSource) {
			if i%2 == 0 {
				continue
			}
			l = append(l, matchSource[i])
		}
		matchSource = l
	}

	if len(matchSource) != len(matchView) {
		err = fmt.Errorf("len(matchSource) != len(matchView), len(matchSource): %d, len(matchView): %d", len(matchSource), len(matchView))
		return
	}

	for i, item := range matchView {
		pr = append(pr, preview.ContentResource{Source: fmt.Sprintf("https://%s", matchSource[i]), View: item})
	}
	return
}

func (p PWorld) Name() string {
	return "pornworld"
}

func New() *PWorld {
	return &PWorld{}
}

// Package allanal do something
package allanal

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type allanal struct {
	Source string
}

func (a allanal) Name() string {
	return "allanal"
}

func (a *allanal) SetSource(source string) {
	a.Source = source
}

func New() *allanal {
	return &allanal{}
}

type JSONAdapter struct {
	Props struct {
		PageProps struct {
			Contents struct {
				Data []struct {
					Title   string `json:"title"`
					Trailer string `json:"trailer_url"`
					Link    string `json:"link"`
				} `json:"data"`
			} `json:"contents"`
		} `json:"pageProps"`
	} `json:"props"`
}

func (allanal) Domain() []string { return Domain }

var Domain = []string{`allanal`, `trueanal`, `analonly`, `swallowed`, `nympho`, `dirtyauditions`}

func (a allanal) Get(index int) (list []preview.ContentResource, err error) {
	uri := fmt.Sprintf(`https://tour.%s.com/scenes`, a.Source)
	if a.Source == `dirtyauditions` {
		uri = "https://dirtyauditions.com/scenes"
	}
	client := resty.New()
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}

	submatch := regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">([^<]+)</script>`).FindStringSubmatch(res.String())
	if len(submatch) < 2 {
		err = fmt.Errorf("len(submatch) < 2")
		return
	}
	var ja JSONAdapter
	err = json.Unmarshal([]byte(submatch[1]), &ja)
	if err != nil {
		return
	}

	for _, i := range ja.Props.PageProps.Contents.Data {
		list = append(list, preview.ContentResource{Source: i.Link, View: i.Trailer})
	}
	return
}

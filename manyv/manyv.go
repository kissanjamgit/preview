// Package manyv is a preview provider for manyvids
package manyv

import (
	"encoding/json"
	"fmt"
	"regexp"

	pr "github.com/kissanjamgit/preview"

	"resty.dev/v3"
)

type preview struct {
	URL string `json:"url"`
}
type resoruce struct {
	ID      string  `json:"id"`
	Preview preview `json:"preview"`
}
type body struct {
	Data []resoruce `json:"data"`
}

type Manyv struct {
	source string
}

func New(source string) *Manyv {
	return &Manyv{source}
}

func (m *Manyv) Search(index int) (list []pr.ContentResource, err error) {
	if index <= 0 {
		index = 1
	}
	var ID string
	if match := regexp.MustCompile("([^/]+)/([0-9]+)").FindStringSubmatch(m.source); len(match) <= 2 {
		err = fmt.Errorf("no match")
		return
	} else {
		ID = match[2]
	}
	client := resty.New()
	res, err := client.R().Get(fmt.Sprintf(`https://api.manyvids.com/store/videos/%s?sort=featured&page=%d`, ID, index))
	if err != nil {
		return
	}
	body := body{}
	json.Unmarshal(res.Bytes(), &body)
	for _, i := range body.Data {
		list = append(list, pr.ContentResource{Source: fmt.Sprintf("https://www.manyvids.com/Video/%s", i.ID), View: i.Preview.URL})
	}
	return
}

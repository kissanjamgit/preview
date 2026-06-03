// Package xp do something
package xp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type xp struct {
	source string
}

func (x xp) Name() string {
	return "xp"
}

func (x *xp) SetSource(source string) {
	x.source = source
}

func (x xp) SearchIdentity() {}

func New() *xp {
	return &xp{}
}

var baseURL = `https://xpxp.eu/`

func (x xp) get(query []string, index int) (list []preview.ContentResource, err error) {
	index += 1
	var uri string
	if len(query) == 0 {
		uri = fmt.Sprintf(`%s/?page=%d`, baseURL, index)
	} else {
		uri = fmt.Sprintf(`%s/tags/%s?page=%d`, baseURL, strings.Join(query, `%20`), index)
	}
	client := resty.New()
	defer client.Close()
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`data-preview="//([^/]+/(\d+)\.mp4)"`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatch {
		if len(m) < 3 {
			continue
		}
		list = append(list, preview.ContentResource{Source: fmt.Sprintf(`%s/videos/%s`, baseURL, m[2]), View: fmt.Sprintf(`https://%s`, m[1])})
	}
	return
}

func (x xp) Get(index int) (list []preview.ContentResource, err error) {
	return x.get(nil, index)
}

func (x xp) Search(query []string, index int) ([]preview.ContentResource, error) {
	return x.get(query, index)
}

// Package xp do something
package xp

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type xp struct {
	source string
}

func New(source string) preview.Preview {
	return &xp{source}
}

var baseURL = `https://pornxp.ph`

func (x *xp) Get(index int) (list []preview.ContentResource, err error) {
	index += 1
	uri := fmt.Sprintf(`%s/?page=%d`, baseURL, index)
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

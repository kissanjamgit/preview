// Package hhshot is some description
package hhshot

import (
	"fmt"
	"os"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type hhshot struct {
	source string
}

func New(source string) preview.Preview {
	return &hhshot{source: source}
}

var baseURL = `https://hookuphotshot.com`

func (p *hhshot) Get(index int) (list []preview.ContentResource, err error) {
	client := resty.New()
	defer client.Close()
	URI := fmt.Sprintf(`%s/nn/categories/movies/%d/latest/`, baseURL, index)
	res, err := client.R().Get(URI)
	os.WriteFile(`content.html`, res.Bytes(), 0o777)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`<a href="([^"]+)"[^>]*>\s*<img[^>]*src0_1x="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatch {
		list = append(list, preview.ContentResource{Source: m[1], View: baseURL + m[2]})
	}
	return
}

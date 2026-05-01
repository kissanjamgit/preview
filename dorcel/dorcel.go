// Package dorcel some
package dorcel

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type dorcel struct {
	source string
}

func New(source string) preview.Preview {
	return &dorcel{source}
}

var (
	baseURL = `https://www.dorcelclub.com`
	headers = map[string]string{
		"Origin":           baseURL,
		`X-Requested-With`: `XMLHttpRequest`,
	}
)

func (d *dorcel) Get(index int) (list []preview.ContentResource, err error) {
	index += 1
	uri := fmt.Sprintf(`%s/scene/list/more/?lang=en&page=%d&sorting=new`, baseURL, index) // 0 isn't allowed; 0 == 1
	client := resty.New()
	res, err := client.R().SetHeaders(headers).Post(uri)
	if err != nil {
		return
	}
	submatchHref := regexp.MustCompile(`href="([^"]+)" class="thumb">\s*<span[^>]*data-hq-preview="([^"]+)">`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatchHref {
		if len(m) < 2 {
			continue
		}
		list = append(list, preview.ContentResource{Source: baseURL + m[1], View: m[2]})

	}

	return
}

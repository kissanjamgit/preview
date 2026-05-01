// Package newsens some
package newsens

import (
	"fmt"
	"os"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type newsens struct {
	Source string
}

func New(source string) preview.Preview {
	return &newsens{source}
}

func (*newsens) Get(index int) (list []preview.ContentResource, err error) {
	index += 1
	uri := fmt.Sprintf(`https://www.newsensations.com/tour_ns/categories/movies_%d_d.html`, index) // 0 isn't allowed;
	client := resty.New()
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	listSrc := []string{}
	listHref := []string{}
	submatchSrc := regexp.MustCompile(`source src="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	os.WriteFile(`content.html`, res.Bytes(), 0o677)
	for _, m := range submatchSrc {
		if len(m) < 2 {
			continue
		}
		listSrc = append(listSrc, m[1])
	}
	submatchHref := regexp.MustCompile(`<h4><a href="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
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

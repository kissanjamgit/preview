// Package czechav some
package czechav

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type czechav struct {
	source string
}

func New(source string) preview.Preview {
	return &czechav{source}
}

var Domain = []string{
	`czechstreets`, `czechfantasy`, `czechgangbang`, `czechmassage`, `perversefamily`,
}

func (c *czechav) Get(index int) (list []preview.ContentResource, err error) {
	baseURL := `https://` + c.source + `.com`
	index += 1
	uri := fmt.Sprintf(`%s/pages/page-%d/?render=1`, baseURL, index)
	client := resty.New()
	defer client.Close()
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	listSrc := []string{}
	listHref := []string{}
	submatch := regexp.MustCompile(`<a class="media-wrapper"  href="([^"]+)">`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatch {
		if len(m) < 2 {
			continue
		}
		listHref = append(listHref, m[1])
	}
	submatch = regexp.MustCompile(`class="poster vthumb"\s*data-thumbnail="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatch {
		if len(m) < 2 {
			continue
		}
		listSrc = append(listSrc, m[1])
	}
	if len(listSrc) != len(listHref) {
		err = fmt.Errorf("len(listSrc) != len(listHref)")
		return
	}
	for i, src := range listSrc {
		list = append(list, preview.ContentResource{Source: baseURL + listHref[i], View: src})
	}
	return
}

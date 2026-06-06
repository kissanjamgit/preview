// Package wh is something
package wh

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/kissanjamgit/preview"
	"github.com/kissanjamgit/preview/common"
	"resty.dev/v3"
)

type wh struct {
	Source string
}

func New() *wh {
	return &wh{
		Source: ``,
	}
}

func (*wh) Name() string {
	return "whoreshub"
}

func (w *wh) SetSource(source string) {
	w.Source = source
}

func (w *wh) get(query []string, index int) (list []preview.ContentResource, err error) {
	common.PlayerArgs = append(common.PlayerArgs, `--http-header-fields=Referer: https://www.whoreshub.com/`)
	client := resty.New()
	uri := "https://www.whoreshub.com/"
	if len(query) != 0 {
		slug := strings.Join(query, `-`)
		uri = fmt.Sprintf("https://www.whoreshub.com/search/%s/", slug)
	}
	q := url.Values{}
	q.Set(`from_videos`, strconv.Itoa(index+1))
	uri += `?` + q.Encode()
	// fmt.Println(uri)
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`alt="([^"]*)"\s*data-preview="([^"]*)"`).FindAllStringSubmatch(res.String(), -1)
	for _, match := range submatch {
		if len(match) < 3 {
			continue
		}
		list = append(list, preview.ContentResource{Source: match[1], View: match[2]})
	}
	return
}

func (w *wh) Get(index int) (list []preview.ContentResource, err error) {
	return w.get(nil, index)
}

// https://www.whoreshub.com/search/piss/
func (*wh) SearchIdentity() {}

func (w *wh) Search(query []string, index int) ([]preview.ContentResource, error) {
	return w.get(query, index)
}

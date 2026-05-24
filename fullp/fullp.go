// Package fullp something something
package fullp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type fullp struct {
	Source string
}

func (f *fullp) SetSource(source string) {
	f.Source = source
}

func New() *fullp {
	return &fullp{}
}

func (f fullp) Name() string {
	return "fullporn"
}

func (f fullp) get(query []string, index int) (list []preview.ContentResource, err error) {
	index += 1
	cliend := resty.New()
	var uri string
	if len(query) != 0 {
		uri = fmt.Sprintf("https://www.fullporn.xxx/search/%s/%d/", strings.Join(query, "-"), index)
	} else {
		uri = fmt.Sprintf("https://www.fullporn.xxx/latest-updates/%d/", index)
	}
	res, err := cliend.R().Get(uri)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`href="([^"]+)"[^>]*>\s*<div[^>]*data-preview="([^"]*)">`).FindAllStringSubmatch(res.String(), -1)
	for _, item := range submatch {
		if len(item) < 3 {
			continue
		}
		list = append(list, preview.ContentResource{Source: item[1], View: item[2]})
	}
	return
}

func (f fullp) Get(index int) (list []preview.ContentResource, err error) {
	return f.get(nil, index)
}
func (f fullp) SearchIdentity() {}
func (f fullp) Search(query []string, index int) (list []preview.ContentResource, err error) {
	return f.get(query, index)
}

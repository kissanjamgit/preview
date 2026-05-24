// Package hqp something something
package hqp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kissanjamgit/preview"
	"github.com/kissanjamgit/preview/common"
	"resty.dev/v3"
)

type hqp struct {
	Source string
}

func (f *hqp) SetSource(source string) {
	f.Source = source
}

func New() *hqp {
	return &hqp{}
}

func (f hqp) Name() string {
	return "hqporner"
}

const baseURL = `https://hqporner.com`

func (f hqp) get(query []string, index int) (list []preview.ContentResource, err error) {
	common.PlayerArgs = append(common.PlayerArgs, `--http-header-fields=Referer: https://hqporner.com/`)
	index += 1
	cliend := resty.New()
	var uri string
	if len(query) != 0 {
		uri = fmt.Sprintf("%s/?q=%s&p=%d", baseURL, strings.Join(query, "+"), index)
	} else {
		uri = fmt.Sprintf("%s/hqporn/%d", baseURL, index)
	}
	res, err := cliend.R().Get(uri)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`href="([^"]+)"[^>]*>\s*<div[^>]*onmouseleave='defaultImage\("([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, item := range submatch {

		if len(item) < 3 {
			continue
		}

		list = append(list, preview.ContentResource{Source: baseURL + item[1], View: `https:` + item[2]})
	}
	return
}

func (f hqp) Get(index int) (list []preview.ContentResource, err error) {
	return f.get(nil, index)
}
func (f hqp) SearchIdentity() {}
func (f hqp) Search(query []string, index int) (list []preview.ContentResource, err error) {
	return f.get(query, index)
}

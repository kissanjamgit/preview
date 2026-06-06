// Package pbk something
package pbk

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/kissanjamgit/preview"
	"github.com/kissanjamgit/preview/config"
	"resty.dev/v3"
)

type pbk struct {
	Source string
}

func New() *pbk {
	return &pbk{}
}

func (p *pbk) SetSource(source string) {
	p.Source = source
}

func (*pbk) Name() string {
	return "premium-bukkake"
}

func (p *pbk) Get(index int) (list []preview.ContentResource, err error) {
	cfg := config.ConfigLazy
	if cfg == nil {
		err = fmt.Errorf("config.ConfigLazy == nil")
		return
	}
	client := resty.New()

	page := index + 1
	uri := `https://premiumbukkake.com/tour2/updates/page_` + strconv.Itoa(page) + `.html`
	res, err := client.SetProxy(config.ConfigLazy.Proxy).R().Get(uri)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`alt="([^"]*) Player Thumbnail">\s*<span[^>]*'([^']+mp4)'`).FindAllStringSubmatch(res.String(), -1)
	for _, match := range submatch {
		list = append(list, preview.ContentResource{Source: match[1], View: match[2]})
	}
	return
}

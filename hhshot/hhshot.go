// Package hhshot is some description
package hhshot

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type hhshot struct {
	source string
}

func (h hhshot) Name() string {
	return "hhshot"
}

func (h *hhshot) SetSource(source string) {
	h.source = source
}

func New() *hhshot {
	return &hhshot{}
}

type siteConfig struct {
	host      string
	route     string
	sourceFmt string
}

var path = []siteConfig{{`hookuphotshot`, `https://hookuphotshot.com/nn/categories/movies/%d/latest/`, `https://hookuphotshot.com%s`}, {`missax`, `https://missax.com/tour/categories/movies_%d_d.html`, `%s`}}

var Domain = func() (list []string) {
	for _, p := range path {
		list = append(list, p.host)
	}
	return
}()

func (p hhshot) Get(index int) (list []preview.ContentResource, err error) {
	config, err := func() (siteConfig, error) {
		for _, r := range path {
			if r.host != p.source {
				continue
			}
			return r, nil
		}
		return siteConfig{}, fmt.Errorf("route not found")
	}()
	if err != nil {
		return
	}
	client := resty.New()
	defer client.Close()
	URI := fmt.Sprintf(config.route, index)
	res, err := client.R().Get(URI)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`<a href="([^"]+)"[^>]*>\s*<img[^>]*src0_1x="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatch {
		list = append(list, preview.ContentResource{Source: m[1], View: fmt.Sprintf(config.sourceFmt, m[2])})
	}
	return
}

// Package pimpbunny something
package pimpbunny

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/kissanjamgit/preview"
	"github.com/kissanjamgit/preview/common"
	"github.com/kissanjamgit/preview/config"
	"resty.dev/v3"
)

type pimpbunny struct {
	Source string
}

func New() *pimpbunny {
	return &pimpbunny{}
}

func (p *pimpbunny) SetSource(source string) {
	p.Source = source
}

func (*pimpbunny) Name() string {
	return "Pimpbunny"
}

const baseURL = `https://pimpbunny.com`

func CookiesToString(cookies []*http.Cookie) string {
	parts := make([]string, 0, len(cookies))

	for _, c := range cookies {
		parts = append(parts, c.Name+"="+c.Value)
	}

	return strings.Join(parts, "; ")
}

func (p *pimpbunny) get(query []string, index int) (list []preview.ContentResource, err error) {
	page := index + 1
	var uri string
	cfg := config.ConfigLazy
	if cfg == nil {
		err = fmt.Errorf("config.ConfigLazy == nil")
		return
	}

	if len(query) != 0 {
		uri = fmt.Sprintf(`%s/search/?q=%s&mode=async&function=get_block&block_id=list_videos_videos_list_search_result&from_videos=%d&ipp=23&page_type=&items_per_page=23&videos_per_page=23&_=1781274060950`, baseURL, strings.Join(query, `%20`), page)
	} else {
		uri = fmt.Sprintf(`%s/videos/%d/`, baseURL, page)
	}
	client := resty.New()

	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	if res.StatusCode() != 200 {
		err = fmt.Errorf("status code %d", res.StatusCode())
		return
	}
	cookieStr := CookiesToString(res.Cookies())
	common.PlayerArgs = append(common.PlayerArgs, fmt.Sprintf(`--http-header-fields=Cookie: %s`, cookieStr))

	submatchSource := regexp.MustCompile(`class="ui-card-link__KxRw6l"\s*href="([^"]*)"`).FindAllStringSubmatch(res.String(), -1)
	submatchView := regexp.MustCompile(`class="\s*ui-card-thumbnail__8dZcLX\s*lazy-load\s*"[^>]*data-preview="([^"]*)"`).FindAllStringSubmatch(res.String(), -1)
	if len(submatchSource) != len(submatchView) {
		err = fmt.Errorf("len(submatchSource) != len(submatchView), len(submatchSource): %d, len(submatchView): %d", len(submatchSource), len(submatchView))
		return
	}

	for i, matchSource := range submatchSource {
		if len(submatchSource) < 2 || len(submatchView) < 2 {
			continue
		}
		list = append(list, preview.ContentResource{Source: matchSource[1], View: submatchView[i][1]})
	}
	return
}

func (p *pimpbunny) Get(index int) (list []preview.ContentResource, err error) {
	return p.get(nil, index)
}

func (p *pimpbunny) SearchIdentity() {}
func (p *pimpbunny) Search(query []string, index int) (list []preview.ContentResource, err error) {
	return p.get(query, index)
}

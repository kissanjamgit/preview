// Package blacked do something
package blacked

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type Blacked struct {
	Source string
}

func (b Blacked) Name() string {
	return "blacked"
}

func (b *Blacked) SetSource(source string) {
	b.Source = source
}

func New() *Blacked {
	return &Blacked{}
}

type JSONAdapte struct {
	Props struct {
		PageProps struct {
			Edges []struct {
				Node struct {
					TITLE string `json:"title"`
					Slug  string `json:"slug"`
					Pre   struct {
						Listing []struct {
							Src string `json:"src"`
						} `json:"listing"`
					} `json:"previews"`
				} `josn:"node"`
			} `json:"edges"`
		} `json:"pageProps"`
	} `json:"props"`
}

var Domain = []string{`blacked`, `tushy`, `vixen`, `slayed`, `milfy`, `wifey`, `blackedraw`, `tushyraw`}

var header = map[string]string{
	`User-Agent`:                `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:150.0) Gecko/20100101 Firefox/150.0`,
	`Accept`:                    `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`,
	`Accept-Language`:           `en-US,en;q=0.9`,
	`Sec-GPC`:                   `1`,
	`Connection`:                `keep-alive`,
	`Upgrade-Insecure-Requests`: `1`,
	`Sec-Fetch-Dest`:            `document`,
	`Sec-Fetch-Mode`:            `navigate`,
	`Sec-Fetch-Site`:            `none`,
	`Sec-Fetch-User`:            `?1`,
	`Priority`:                  `u=0, i`,
}

func (b Blacked) Get(index int) (list []preview.ContentResource, err error) {
	index += 1
	baseURL := fmt.Sprintf(`https://www.%s.com`, b.Source)
	uri := fmt.Sprintf(`%s/videos?page=%d`, baseURL, index)
	client := resty.New()
	res, err := client.R().SetHeaders(header).Get(uri)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`<script id="__NEXT_DATA__"[^>]*>([^<]+)</script`).FindStringSubmatch(res.String())
	if len(submatch) < 2 {
		err = fmt.Errorf("len(submatch) < 2")
		return
	}
	var ja JSONAdapte
	json.Unmarshal([]byte(submatch[1]), &ja)
	for _, item := range ja.Props.PageProps.Edges {
		viewList := item.Node.Pre.Listing
		view := viewList[len(viewList)-1].Src
		list = append(list, preview.ContentResource{Source: fmt.Sprintf(`%s/videos/%s`, baseURL, item.Node.Slug), View: view})
	}
	return
}

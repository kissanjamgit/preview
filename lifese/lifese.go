// Package lifese
package lifese

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/ext"
	"resty.dev/v3"
)

type Lifese struct {
	source string
}

func New(source string) ext.Site {
	return &Lifese{source: source}
}

var header = map[string]string{
	"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Language":           "en-US,en;q=0.5",
	"Referer":                   "https://lifeselector.com/",
	"Sec-GPC":                   "1",
	"Connection":                "keep-alive",
	"Cookie":                    "SID=7n88essp49enaqqtl5ru0k32n0; locale=en.1799077570; __cflb=0H28vTz59LtvvwE3hA1BLesBN9zumY2wbGis57oFhyR; age_verification=1",
	"Upgrade-Insecure-Requests": "1",
	"Sec-Fetch-Dest":            "document",
	"Sec-Fetch-Mode":            "navigate",
	"Sec-Fetch-Site":            "same-origin",
	"Sec-Fetch-User":            "?1",
	"Priority":                  "u=0, i",
	"TE":                        "trailers",
}

func (s *Lifese) Resource(client *resty.Client) (cr ext.ContentResource, err error) {
	submatchID := regexp.MustCompile(`gameId/(\d+)`).FindStringSubmatch(s.source)
	if len(submatchID) < 1 {
		return
	}
	URI := fmt.Sprintf(`https://lifeselector.com/game/trailer/gameId/%s/ext/file.mp4`, submatchID[1])
	res, err := client.R().SetHeaders(header).Get(s.source)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`<h1 class="title">([^<]+)</h1>`).FindStringSubmatch(res.String())
	if len(submatch) < 1 {
		err = fmt.Errorf("submatch len < 1")
		return
	}
	cr = ext.ContentResource{Name: submatch[1], URL: URI}
	return
}

func (s *Lifese) Download(cr ext.ContentResource) (err error) {
	return
}

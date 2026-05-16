// Package lifese do something
package lifese

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type lifese struct {
	source string
}

func (life lifese) Name() string {
	return "lifese"
}

func (life *lifese) SetSource(source string) {
	life.source = source
}

func New() *lifese {
	return &lifese{}
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

var size = 21

func (life lifese) Get(index int) (list []preview.ContentResource, err error) {
	URI := url.URL{
		Scheme: "https",
		Host:   "lifeselector.com",
		Path:   "game/listgames",
	}
	q := url.Values{}
	q.Add("format", "partial")
	q.Add("offset", fmt.Sprintf("%d", size*index))
	q.Add("order", "releaseDate")
	q.Add("_", strconv.Itoa(int(time.Now().UnixMilli())))
	URI.RawQuery = q.Encode()

	client := resty.New()
	defer resty.New().Close()

	res, err := client.R().SetHeaders(header).Get(URI.String())
	if err != nil {
		return
	}

	submatch := regexp.MustCompile(`class="cover  game-4k" href="([^"]+)">\s*<img src="([^"]+)"`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatch {
		if len(m) < 3 {
			continue
		}
		list = append(list, preview.ContentResource{Source: fmt.Sprintf(`https://%s%s`, URI.Host, m[1]), View: m[2]})
	}

	return
}

// Package lifese
package lifese

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/kissanjamgit/ext"
	"golang.org/x/sync/errgroup"
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

type jsonAdapte struct {
	Resource struct {
		Content string `json:"content"`
	} `json:"resource"`
}

// Resource Note cr.URL for lifese is a redirect it requres to have header "TE:trailers"
// add this for mpvnet to work with it "ytdl-raw-options=extractor-args=add-headers=TE:trailers" //this shorter form of the value which is not been tested

func (s *Lifese) Resource(client *resty.Client) (cr ext.ContentResource, err error) {
	submatchID := regexp.MustCompile(`gameId/(\d+)`).FindStringSubmatch(s.source)
	if len(submatchID) < 2 {
		return
	}

	var name string
	g := errgroup.Group{}
	g.Go(func() (err error) {
		res, err := client.R().SetHeaders(header).Get(s.source)
		if err != nil {
			return
		}
		submatch := regexp.MustCompile(`<h1 class="title">([^<]+)</h1>`).FindStringSubmatch(res.String())
		if len(submatch) < 2 {
			err = fmt.Errorf("len(submatch) < 2 ")
			return
		}
		name = submatch[1]
		return
	})

	URI := fmt.Sprintf(`https://lifeselector.com/game/GetEpisodeDetailsInJson/gameId/%s/choiceId/0`, submatchID[1])
	res, err := client.R().SetHeaders(header).Get(URI)
	if err != nil {
		return
	}
	var ja jsonAdapte
	err = json.Unmarshal(res.Bytes(), &ja)
	if err != nil {
		return
	}
	err = g.Wait()
	if err != nil {
		return
	}
	cr = ext.ContentResource{Name: name, URL: ja.Resource.Content}
	return
}

func (s *Lifese) Download(cr ext.ContentResource) (err error) {
	cmd := exec.Command("yt-dlp.exe", cr.URL, "-o", cr.Name+".mp4", "--add-header", "TE:trailers", "--extractor-args", "generic:impersonate")
	err = cmd.Run()
	return
}

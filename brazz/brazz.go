// Package brazz something something
package brazz

import (
	"encoding/json"
	"fmt"
	"maps"
	"regexp"
	"strings"
	"time"

	"github.com/kissanjamgit/preview"

	"resty.dev/v3"
)

var Domain = []string{
	"brazzers",
	"milfed",
	"bangbros",
	"realitykings",
	"twistys",
}

type Brazz struct {
	Source string
}

func (b *Brazz) Get(index int) (list []preview.ContentResource, err error) {
	client := resty.New()
	defer client.Close()
	var jwt string

	domain := fmt.Sprintf("https://%s.com/", b.Source)
	req, err := client.R().SetHeaders(header(nil, domain)).Get(domain)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := regexp.MustCompile(`"jwt":"(\S*?)"`).FindStringSubmatch(req.String())
	jwt = data[len(data)-1]

	url := fmt.Sprintf("https://site-api.project1service.com/v2/releases?adaptiveStreamingOnly=false&dateReleased=%%3C%s&orderBy=-dateReleased&type=scene&limit=40&offset=%d", time.Now().Format("2006-01-02"), 40*index)
	res, err := client.R().SetHeaders(header(&jwt, domain)).Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	var adapt AdaptarJSONParse
	err = json.Unmarshal(res.Bytes(), &adapt)
	if err != nil {
		return
	}

	for _, item := range adapt.Result {
		list = append(list, preview.ContentResource{Source: cleanBrazz(b.Source, item), View: item.Video.Mediabook.Files.Res720p.Urls.View})
	}
	return
}

func New(source string) preview.Preview {
	return &Brazz{source}
}

var sceneOrVideo = map[string]bool{
	"milfed":            true,
	"mofos":             true,
	"twistys":           true,
	"realitykings":      true,
	"digitalplayground": true,
}

var (
	invalidChars = regexp.MustCompile(`[^a-z0-9\s]`)
	header_      = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0",
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "en-US,en;q=0.5",
		"Origin":          "https://www.brazzers.com",
		"Sec-GPC":         "1",
		"Connection":      "keep-alive",
		"Referer":         "https://www.brazzers.com/",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "cross-site",
		"TE":              "trailers",
	}
)

func header(jwt *string, domain string) map[string]string {
	h := map[string]string{}
	maps.Copy(h, header_)
	if jwt != nil && *jwt != "" {
		h["instance"] = *jwt
	}
	h["Origin"] = domain
	h["Referer"] = domain
	return h
}

func cleanBrazz(domain string, j AdapatJSONResult) string {
	s := j.Title
	s = strings.ToLower(s)
	s = invalidChars.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, " ", "-")

	S := "scene"
	if b := sceneOrVideo[domain]; !b {
		S = "video"
	}
	return fmt.Sprintf("https://%s.com/%s/%v/%s", domain, S, j.ID, s)
}

type AdapatJSONVideo struct {
	Mediabook struct {
		Files struct {
			Res720p struct {
				Urls struct {
					View string `json:"view"`
				} `json:"urls"`
			} `json:"720p"`
		} `json:"files"`
	} `json:"mediabook"`
}

type AdapatJSONResult struct {
	ID    int             `json:"id"`
	Title string          `json:"title"`
	Video AdapatJSONVideo `json:"videos"`
}

type AdaptarJSONParse struct {
	Result []AdapatJSONResult `json:"result"`
}

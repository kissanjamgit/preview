package hooporn

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type hooporn struct {
	Source string
}

func New() *hooporn {
	return &hooporn{}
}

func (h *hooporn) SetSource(source string) {
	h.Source = source
}

func (h *hooporn) Name() string {
	return "hooporn"
}

func (h *hooporn) Get(index int) ([]preview.ContentResource, error) {
	client := resty.New()
	client.SetHeaders(map[string]string{
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:152.0) Gecko/20100101 Firefox/152.0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language":           "en-US,en;q=0.9",
		"Accept-Encoding":           "gzip, deflate, br, zstd",
		"Sec-GPC":                   "1",
		"Connection":                "keep-alive",
		"Referer":                   "https://hooporn.com/",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "same-origin",
		"Sec-Fetch-User":            "?1",
		"Priority":                  "u=0, i",
		"Pragma":                    "no-cache",
		"Cache-Control":             "no-cache",
		"TE":                        "trailers",
	})

	page := index + 1
	url := fmt.Sprintf("https://hooporn.com/home/%d", page)
	res, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	// Regex to match the video links
	re := regexp.MustCompile(`href="(/watch/[^"]+)"`)
	matches := re.FindAllStringSubmatch(res.String(), -1)

	list := make([]preview.ContentResource, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		fullURL := "https://hooporn.com" + match[1]
		list = append(list, preview.ContentResource{
			Source: fullURL,
			View:   fullURL,
		})
	}

	return list, nil
}

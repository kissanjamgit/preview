package kink

import (
	"fmt"
	"os"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type kink struct {
	Source string
}

func New() *kink {
	return &kink{}
}

func (k *kink) SetSource(source string) {
	k.Source = source
}

func (k *kink) Name() string {
	return "kink"
}

func (k *kink) Get(index int) ([]preview.ContentResource, error) {
	client := resty.New()
	client.SetHeaders(map[string]string{
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:152.0) Gecko/20100101 Firefox/152.0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language":           "en-US,en;q=0.9",
		"Referer":                   "https://www.kink.com/",
		"Sec-GPC":                   "1",
		"Connection":                "keep-alive",
		"Cookie":                    "deviceId=ced52e36-b762-41b3-99c7-b349cd473e52; ct=2; kinky.sess=s%3Ac8b311c0-8025-11f1-99f0-3bedf76e66b5.YwXHhFeEQP87rWp1iJ0h2A9veaFwv9T8IJQCsVxtvFc; age_gate_accepted=1",
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

	res, err := client.R().Get("https://www.kink.com/shoots?sort=published")
	if err != nil {
		return nil, err
	}

	// Regex to capture href using the specified class
	os.WriteFile("context.html", []byte(res.String()), 0o644)
	re := regexp.MustCompile(`class="d-block overflow-hidden text-elipsis h5"[^>]*href="([^"]*)"`)
	matches := re.FindAllStringSubmatch(res.String(), -1)

	list := make([]preview.ContentResource, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		// match[1] is the relative URL
		fullURL := fmt.Sprintf("https://www.kink.com%s", match[1])

		list = append(list, preview.ContentResource{
			Source: fullURL,
			View:   fullURL,
		})
	}
	return list, nil
}

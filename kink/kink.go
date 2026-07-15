// Package kink provides a kink.com scraper
package kink

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

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

func (k *kink) SearchIdentity() {}

func (k *kink) Search(query []string, index int) ([]preview.ContentResource, error) {
	q := strings.Join(query, " ")
	uri := fmt.Sprintf("https://www.kink.com/shoots?sort=published&q=%s", url.QueryEscape(q))
	return k.fetch(uri)
}

func (k *kink) Get(index int) ([]preview.ContentResource, error) {
	uri := fmt.Sprintf("https://www.kink.com/shoots?sort=published&page=%d", index+1)
	return k.fetch(uri)
}

func (k *kink) fetch(uri string) ([]preview.ContentResource, error) {
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

	res, err := client.R().Get(uri)
	if err != nil {
		return nil, err
	}
	// os.WriteFile("result.txt", res.Bytes(), 0o644)

	regexpTrailer := regexp.MustCompile(`<img[^>]*((data-trailer-url="[^"]+")|(data-src="[^"]+"))[^>]*class="has-kink-spinner"\/>`).FindAllStringSubmatch(res.String(), -1) //<img[^>]*(data-trailer-url="([^"]+)")?[^>]*(data-cycle="([^"]+)")?

	list := make([]preview.ContentResource, 0, len(regexpTrailer))

	regexpShootID := regexp.MustCompile(`shoots/(\d+)/|imagedb/(\d+)/`)
	for _, match := range regexpTrailer {

		var shootID string
		shootIDMatch := regexpShootID.FindStringSubmatch(match[0])
		if len(shootIDMatch) < 2 {
			continue
		}
		if strings.HasPrefix(shootIDMatch[0], `shoots`) {
			shootID = shootIDMatch[1]
		} else {
			shootID = shootIDMatch[2]
		}
		source := "https://www.kink.com/shoot/" + shootID

		var view string
		if after, ok := strings.CutPrefix(match[2], `data-trailer-url=`); ok {
			view = strings.Trim(after, `"`)
		} else if after, ok := strings.CutPrefix(match[3], `data-src=`); ok {
			view = strings.Trim(after, `"`)
			// view = regexpHTTP.FindString(after)
			// view = strings.TrimSuffix(view, `&quot`)
		} else {
			continue
			// return nil, fmt.Errorf(`prefix doesn't match %s`, match[2])
		}

		list = append(list, preview.ContentResource{
			Source: source,
			// Source: `a`,
			View: view,
		})
	}
	return list, nil
}

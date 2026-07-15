// pacakge netfapx provides preview functionality for netfapx.com
package netfapx

import (
	"os"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type netfapx struct {
	Source string
}

func New() *netfapx {
	return &netfapx{}
}

func (n *netfapx) SetSource(source string) {
	n.Source = source
}

func (n *netfapx) Name() string {
	return "netfapx"
}

func (n *netfapx) Get(index int) ([]preview.ContentResource, error) {
	client := resty.New()
	client.SetHeaders(map[string]string{
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:152.0) Gecko/20100101 Firefox/152.0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language":           "en-US,en;q=0.9",
		"Accept-Encoding":           "gzip, deflate, br, zstd",
		"Sec-GPC":                   "1",
		"Connection":                "keep-alive",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Priority":                  "u=4",
		"Pragma":                    "no-cache",
		"Cache-Control":             "no-cache",
	})

	res, err := client.R().Get("https://netfapx.com/")
	if err != nil {
		return nil, err
	}
	os.WriteFile(`context.html`, res.Bytes(), 0o644)

	// TODO: Implement extraction logic here once the site is accessible
	return []preview.ContentResource{}, nil
}

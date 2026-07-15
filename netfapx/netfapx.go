// pacakge netfapx provides preview functionality for netfapx.com
package netfapx

import (
	"github.com/go-resty/resty/v2"
	"github.com/kissanjamgit/preview"
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
	// Using a standard Firefox User-Agent to bypass potential blocks
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0")

	_, err := client.R().Get("https://netfapx.com/")
	if err != nil {
		return nil, err
	}

	// TODO: Implement extraction logic here once the site is accessible
	return []preview.ContentResource{}, nil
}

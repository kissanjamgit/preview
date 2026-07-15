// pacakge netfapx provides preview functionality for netfapx.com
package netfapx

import (
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
	return []preview.ContentResource{}, nil
}

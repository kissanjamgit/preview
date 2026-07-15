package porndd

import (
	"github.com/kissanjamgit/preview"
)

type porndd struct {
	Source string
}

func New() *porndd {
	return &porndd{}
}

func (p *porndd) SetSource(source string) {
	p.Source = source
}

func (p *porndd) Name() string {
	return "porndd"
}

func (p *porndd) Get(index int) ([]preview.ContentResource, error) {
	// Placeholder for extraction logic
	return []preview.ContentResource{}, nil
}

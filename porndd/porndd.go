package porndd

import (
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
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
	client := resty.New()
	res, err := client.R().Get("https://porndd.com/")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`<a href="([^"]*)">\s*<img [^>]*src="([^"]*)"`)
	matches := re.FindAllStringSubmatch(res.String(), -1)

	list := make([]preview.ContentResource, 0, len(matches))
	for _, match := range matches {
		list = append(list, preview.ContentResource{
			Source: match[1], // Using href as source
			View:   match[1], // Using href as view link
		})
	}

	return list, nil
}

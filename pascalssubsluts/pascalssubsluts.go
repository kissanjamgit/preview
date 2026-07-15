package pascalssubsluts

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type pascalssubsluts struct {
	Source string
}

var baseURL = "https://www.pascalssubsluts.com"

func New() *pascalssubsluts {
	return &pascalssubsluts{}
}

func (p *pascalssubsluts) SetSource(source string) {
	p.Source = source
}

func (p *pascalssubsluts) Name() string {
	return "pascalssubsluts"
}

func (p *pascalssubsluts) Get(index int) ([]preview.ContentResource, error) {
	client := resty.New()
	// Assuming the updates page is the source of the list
	uri := fmt.Sprintf("%s/submissive/updates.php", baseURL)
	res, err := client.R().Get(uri)
	if err != nil {
		return nil, err
	}

	// Regex to capture data-joinimg and href
	re := regexp.MustCompile(`data-joinimg="([^"]*)"[^>]*href="player-load.php?id=(/d+)"`)
	matches := re.FindAllStringSubmatch(res.String(), -1)

	list := make([]preview.ContentResource, 0, len(matches))
	fmt.Println(matches)
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		// Decode the URL-encoded image source

		list = append(list, preview.ContentResource{
			Source: match[2], // Using href as source
			View:   match[1], // Using decoded image URL as view
		})
	}

	return list, nil
}

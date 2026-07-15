package pascalssubsluts

import (
	"net/url"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type pascalssubsluts struct {
	Source string
}

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
	res, err := client.R().Get("https://www.pascalssubsluts.com/submissive/updates.php")
	if err != nil {
		return nil, err
	}

	// Regex to capture data-joinimg and data-modelname
	re := regexp.MustCompile(`data-joinimg="([^"]*)"[^>]*data-modelname="([^"]*)"`)
	matches := re.FindAllStringSubmatch(res.String(), -1)

	list := make([]preview.ContentResource, 0, len(matches))
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		// Decode the URL-encoded image source
		imgURL, err := url.QueryUnescape(match[1])
		if err != nil {
			continue
		}

		list = append(list, preview.ContentResource{
			Source: match[2], // Using model name as source
			View:   imgURL,   // Using decoded image URL as view
		})
	}

	return list, nil
}

package porndd

import (
	"os"
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
	client.SetHeaders(map[string]string{
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:152.0) Gecko/20100101 Firefox/152.0",
		"Accept":                    "*/*",
		"Accept-Language":           "en-US,en;q=0.9",
		"Referer":                   "https://porndd.com/",
		"Sec-GPC":                   "1",
		"Connection":                "keep-alive",
		"Sec-Fetch-Dest":            "empty",
		"Sec-Fetch-Mode":            "no-cors",
		"Sec-Fetch-Site":            "same-origin",
		"Priority":                  "u=4",
		"Pragma":                    "no-cache",
		"Cache-Control":             "no-cache",
		"TE":                        "trailers",
	})

	res, err := client.R().Get("https://porndd.com/?mode=async&function=get_block&block_id=list_videos_most_recent_videos&sort_by=post_date&from=1/")
	if err != nil {
		return nil, err
	}
	os.WriteFile(`context.html`, res.Bytes(), 0o644)
	
	// Updated regex to handle nested structure and data-original attribute
	re := regexp.MustCompile(`href="([^"]*)"[^>]*>[\s\S]*?data-original="([^"]*)"`)
	matches := re.FindAllStringSubmatch(res.String(), -1)

	list := make([]preview.ContentResource, 0, len(matches))
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		list = append(list, preview.ContentResource{
			Source: match[1], // Using href as source
			View:   match[1], // Using href as view link
		})
	}

	return list, nil
}

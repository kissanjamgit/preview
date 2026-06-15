// Package faphouse provides preview functionality for faphouse.com
package faphouse

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/kissanjamgit/preview"
)

// Faphouse struct for faphouse package
type Faphouse struct {
	Source string // URL domain (e.g., "faphouse.com")
}


// SetSource sets the source URL domain for faphouse
func (f *Faphouse) SetSource(source string) {
	f.Source = source
}

// Name returns the name of this preview provider
func (f *Faphouse) Name() string {
	return "faphouse"
}

// Get returns content resources from the faphouse preview at the given index
func (f *Faphouse) Get(index int) ([]preview.ContentResource, error) {
	if index < 0 {
		return nil, fmt.Errorf("index out of range")
	}

	if f.Source == "" {
		return nil, fmt.Errorf("source is required")
	}

	html, err := f.fetchHTML(f.Source)
	if err != nil {
		return nil, err
	}

	resources := f.parseHTML(html)

	if index >= len(resources) {
		return nil, fmt.Errorf("index out of range")
	}

	var result []preview.ContentResource
	for i := index; i < len(resources); i++ {
		result = append(result, resources[i])
	}

	return result, nil
}


// fetchHTML fetches HTML content from the faphouse source URL domain using resty client
func (f *Faphouse) fetchHTML() ([]preview.ContentResource, error) {
	if f.Source == "" {
		return nil, fmt.Errorf("source is required")
	}

	client := resty.New()
	req, err := client.R().SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	res, err := client.R().Get(f.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %v", err)
	}

	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("status code %d", res.StatusCode())
	}

	html := res.String()

	return f.parseHTML(html), nil
}


// parseHTML parses HTML content and extracts content resources
func (f *Faphouse) parseHTML(html string) []preview.ContentResource {
	var resources []preview.ContentResource

	// Regex to extract data-el-video attribute (source) - uses class names like pimpbunny
	sourceRegex := regexp.MustCompile(`class="thumb_col4 thumb tv"\s*data-el-video="([^"]+)"`)

	// Regex to extract img src attribute (view) - uses class names like pimpbunny
	viewRegex := regexp.MustCompile(`class="t-i"\s*src="([^"]+)"`)

	// Find all source matches
	sourceMatches := sourceRegex.FindAllStringSubmatch(html, -1)

	// Find all view matches
	viewMatches := viewRegex.FindAllStringSubmatch(html, -1)

	// Pair them up (skip if counts don't match or indices are invalid)
	for i := 0; i < len(sourceMatches); i++ {
		if len(sourceMatches[i]) < 2 || len(viewMatches) == 0 {
			continue
		}

		// Find corresponding view for this source
		for j := 0; j < len(viewMatches); j++ {
			if len(viewMatches[j]) >= 2 {
				resources = append(resources, preview.ContentResource{
					Source: sourceMatches[i][1],
					View:   viewMatches[j][1],
				})
				break
			}
		}

		if len(viewMatches) > 0 {
			break
		}
	}

	return resources
}

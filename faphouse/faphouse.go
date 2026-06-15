// Package faphouse provides preview functionality for faphouse.com
package faphouse

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

// Faphouse struct for faphouse package
type Faphouse struct {
	Source string // URL domain (e.g., "faphouse.com")
}

// New creates a new instance of Faphouse with the default domain
func New() *Faphouse {
	return &Faphouse{
		Source: "faphouse.com",
	}
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

	client := resty.New()
	defer client.Close()

	res, err := client.R().Get(f.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %v", err)
	}

	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("status code %d", res.StatusCode())
	}

	html := res.String()

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

	if index >= len(resources) {
		return nil, fmt.Errorf("index out of range")
	}

	var result []preview.ContentResource
	for i := index; i < len(resources); i++ {
		result = append(result, resources[i])
	}

	return result, nil
}

// Domain returns the domain of this preview provider
func (f *Faphouse) Domain() []string {
	return []string{f.Source}
}

// DefaultDomain returns the default domain for this preview provider
func (f *Faphouse) DefaultDomain() string {
	return "faphouse.com"
}

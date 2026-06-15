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
func (f *Faphouse) Get(index int) (list []preview.ContentResource, err error) {
	uri := fmt.Sprintf(`https://faphouse.com/videos?type=new&page=%d`, index+1)
	if index < 0 {
		err = fmt.Errorf("index out of range")
		return
	}

	if f.Source == "" {
		err = fmt.Errorf("source is required")
		return
	}

	client := resty.New()
	defer client.Close()

	res, err := client.R().Get(uri)
	if err != nil {
		err = fmt.Errorf("failed to fetch: %v", err)
		return
	}

	if res.StatusCode() != 200 {
		err = fmt.Errorf("status code %d", res.StatusCode())
		return
	}

	html := res.String()

	var resources []preview.ContentResource

	sourceRegex := regexp.MustCompile(`class="thumb_col4 thumb tv"\s*data-el-video="([^"]+)"`)

	viewRegex := regexp.MustCompile(`class="t-i"\s*src="([^"]+)"`)

	sourceMatches := sourceRegex.FindAllStringSubmatch(html, -1)

	viewMatches := viewRegex.FindAllStringSubmatch(html, -1)

	if len(sourceMatches) == len(viewMatches) {
		err = fmt.Errorf(`len(sourceMatches) == len(viewMatches) `)
		return
	}
	for i, source := range sourceMatches {
		if len(source) < 2 || len(viewMatches) < 2 {
			continue
		}
		resources = append(resources, preview.ContentResource{
			Source: sourceMatches[i][1],
			View:   viewMatches[i][1],
		})

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

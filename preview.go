// Package preview something something
package preview

import (
	"faphouse"
)

type ContentResource struct {
	Source string `json:"name"`
	View   string `json:"view"`
}

type Preview interface {
	SetSource(source string)
	Name() string
	Get(index int) ([]ContentResource, error)
}

type Domain interface {
	Preview
	Domain() []string
}

type Search interface {
	Preview
	SearchIdentity()
	Search(query []string, index int) ([]ContentResource, error)
}

type SearchDomain interface {
	Preview
	SearchIdentity()
	SearchDomainIdentity()
	Domain() []string
	Search(query []string, index int) ([]ContentResource, error)
}

type FaphousePreview struct {
	URL string
	Pairs []struct {
		View  string `json:"view"`
		Source string `json:"source"`
	}
}

func (f *FaphousePreview) SetSource(source string) {
	f.URL = source
}

func (f *FaphousePreview) Name() string {
	return "faphouse"
}

func (f *FaphousePreview) Get(index int) ([]ContentResource, error) {
	if index < 0 || index >= len(f.Pairs) {
		return nil, nil
	}

	var resources []ContentResource
	for i := index; i < len(f.Pairs); i++ {
		resources = append(resources, ContentResource{
			Source: f.Pairs[i].Source,
			View:   f.Pairs[i].View,
		})
	}

	return resources, nil
}

func (f *FaphousePreview) Search(query []string, index int) ([]ContentResource, error) {
	// Implement search logic for faphouse
	return nil, nil
}

func (f *FaphousePreview) SearchIdentity() {
	// Implement search identity for faphouse
}

func (f *FaphousePreview) SearchDomainIdentity() {
	// Implement search domain identity for faphouse
}

func (f *FaphousePreview) Domain() []string {
	return []string{"faphouse"}
}

type Faphouse struct {
	URL string
	Pairs []struct {
		View  string `json:"view"`
		Source string `json:"source"`
	}
	StudioPath faphouse.StudioPath
}

func (f *Faphouse) Fetch() ([]ContentResource, error) {
	pairs, err := f.PairsFromFaphouse()
	if err != nil {
		return nil, err
	}

	var resources []ContentResource
	for _, pair := range pairs {
		resources = append(resources, ContentResource{
			Source: pair.Source,
			View:   pair.View,
		})
	}

	return resources, nil
}

func (f *Faphouse) ParseViewSourcePairs() ([]ContentResource, error) {
	pairs, err := f.PairsFromFaphouse()
	if err != nil {
		return nil, err
	}

	var resources []ContentResource
	for _, pair := range pairs {
		resources = append(resources, ContentResource{
			Source: pair.Source,
			View:   pair.View,
		})
	}

	return resources, nil
}

func (f *Faphouse) PairsFromFaphouse() ([]faphouse.FaphousePair, error) {
	if f.URL == "" {
		return nil, fmt.Errorf("URL is required")
	}

	faphouseClient := faphouse.NewFaphouse(f.URL, "faphouse")
	pairs, err := faphouseClient.FetchHTML()
	if err != nil {
		return nil, err
	}

	return pairs, nil
}

func (f *Faphouse) GetName() string {
	if f.StudioPath.name != "" {
		return f.StudioPath.name
	}
	return "faphouse"
}

func (f *Faphouse) GetPath() string {
	if f.StudioPath.path != "" {
		return f.StudioPath.path
	}
	return "faphouse"
}

func main() {
	// Main entry point for the application
}

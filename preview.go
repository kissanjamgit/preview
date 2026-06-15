// Package preview defines the core interfaces for content previews
package preview

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

// Package preview something something
package preview

type ContentResource struct {
	Source string `json:"name"`
	View   string `json:"view"`
}

type Preview interface {
	Get(index int) ([]ContentResource, error)
}

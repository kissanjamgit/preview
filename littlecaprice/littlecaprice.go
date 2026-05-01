// Package littlecaprice some
package littlecaprice

import (
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type littlecaprice struct {
	source string
}

func New(source string) preview.Preview {
	return &littlecaprice{source}
}

var size = 30

func (l *littlecaprice) Get(index int) (list []preview.ContentResource, err error) {
	index += 1
	uri := `https://www.littlecaprice-dreams.com/videos`
	client := resty.New()
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	submatch := regexp.MustCompile(`href='([^']+)' onmouseenter="LCD_ProjectTrailerPreview\('([^']+)'`).FindAllStringSubmatch(res.String(), -1)
	for _, m := range submatch {
		list = append(list, preview.ContentResource{Source: m[1], View: m[2]})
	}
	list = list[index*size : index*size+size]

	return
}

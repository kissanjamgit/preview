// Package vip
package vip

import (
	"fmt"
	"os"
	"regexp"

	"github.com/kissanjamgit/preview"

	"resty.dev/v3"
)

type Vip struct {
	source string
}

func (v *Vip) SetSource(source string) {
	v.source = source
}

func (v Vip) Get(index int) (pr []preview.ContentResource, err error) {
	if index < 1 {
		index = 1
	}
	url := fmt.Sprintf("https://vip4k.com/en/publish/tag/all/all/all/%d?exclude_sets=1", index)
	client := resty.New()
	res, err := client.R().Get(url)
	if err != nil {
		return
	}
	os.WriteFile("2.html", res.Bytes(), 0o644)
	hrefList := regexp.MustCompile(`item__main"\s+href="([^"]+)`).FindAllStringSubmatch(res.String(), -1)

	prList := regexp.MustCompile(`"//([^"]+preview.mp4)"`).FindAllStringSubmatch(res.String(), -1)

	if len(hrefList) != len(prList) {
		err = fmt.Errorf("len(hrefList) != len(prList)  len(hrefList): %d  len(prList): %d", len(hrefList), len(prList))
		return
	}

	for i, h := range hrefList {
		if len(h) < 2 || len(prList[i]) < 2 {
			err = fmt.Errorf("len(h) < 2 || len(prList[i]) < 2")
			return
		}
		pr = append(pr, preview.ContentResource{Source: fmt.Sprintf("https://vip4k.com%s", h[1]), View: "https://" + prList[i][1]})
	}

	return
}

func (v Vip) Name() string {
	return "vip4k"
}

func New() *Vip {
	return &Vip{}
}

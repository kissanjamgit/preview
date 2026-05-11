// Package advd something
package advd

import (
	"fmt"
	"regexp"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type advd struct {
	source string
}

func New(source string) preview.Preview {
	return &advd{source: source}
}

var (
	baseURL = `https://www.adultdvdempire.com`
	viewFmt = `https://video.adultempire.com/hls/previewscene/%s/%s/index-f2-v1-a1.m3u8`
)

func (a *advd) Get(index int) (list []preview.ContentResource, err error) {
	index += 1
	uri := fmt.Sprintf(`%s/clips/?page=%d`, baseURL, index)
	client := resty.New()
	res, err := client.R().SetHeader("Cookie", "ageConfirmed=true; defaults={}").Get(uri)
	if err != nil {
		return
	}
	submatchData := regexp.MustCompile(`data-scene-id="(\d+)"\s*data-movie-id="(\d+)"`).FindAllStringSubmatch(res.String(), -1)
	submatchHref := regexp.MustCompile(`href='([^"]+)' data-ta="view" `).FindAllStringSubmatch(res.String(), -1)
	if len(submatchData) != len(submatchHref) {
		err = fmt.Errorf(`len(submatchData) != len(submatchHref)`)
		return
	}
	for index, v := range submatchData {
		if len(v) < 3 {
			continue
		}
		matchHref := submatchHref[index]
		if len(matchHref) < 2 {
			continue
		}
		view := fmt.Sprintf(viewFmt, v[2], v[1])
		list = append(list, preview.ContentResource{Source: baseURL + matchHref[1], View: view})
	}

	return
}

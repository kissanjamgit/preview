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

var Domain = []string{`movie`, `clip`}
var (
	baseURL  = `https://www.adultdvdempire.com`
	movieFmt = `https://video.adultempire.com/hls/trailer/%s/index-f2-v1-a1.m3u8`
	clipFmt  = `https://video.adultempire.com/hls/previewscene/%s/%s/index-f2-v1-a1.m3u8`
)

func (a *advd) Get(index int) (list []preview.ContentResource, err error) {
	index += 1
	client := resty.New()
	var extact func(string) ([]preview.ContentResource, error)
	var uri string
	defer client.Close()
	switch a.source {
	case `clip`:
		uri = fmt.Sprintf(`%s/clips/?sort=released&page=%d`, baseURL, index)

		extact = clip
	case `movie`:
		uri = fmt.Sprintf(`%s/new-release-porn-movies.html?page=%d`, baseURL, index)
		extact = movie
	default:
		err = fmt.Errorf("switch case exhausted")
		return
	}
	res, err := client.R().SetHeader("Cookie", "ageConfirmed=true; defaults={}").Get(uri)
	if err != nil {
		return
	}
	list, err = extact(res.String())
	return
}

func movie(s string) (list []preview.ContentResource, err error) {
	submatch := regexp.MustCompile(`href="\s*(/(\d+)[^"]*porn-movies.html)"`).FindAllStringSubmatch(s, -1)
	for _, m := range submatch {
		if len(m) < 2 {
			continue
		}
		list = append(list, preview.ContentResource{Source: baseURL + m[1], View: fmt.Sprintf(movieFmt, m[2])})
	}

	return
}

func clip(s string) (list []preview.ContentResource, err error) {
	submatchData := regexp.MustCompile(`data-scene-id="(\d+)"\s*data-movie-id="(\d+)"`).FindAllStringSubmatch(s, -1)
	submatchHref := regexp.MustCompile(`href='([^"]+)' data-ta="view" `).FindAllStringSubmatch(s, -1)
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
		view := fmt.Sprintf(clipFmt, v[2], v[1])
		list = append(list, preview.ContentResource{Source: baseURL + matchHref[1], View: view})
	}
	return
}

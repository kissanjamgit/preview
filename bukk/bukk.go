package bukk

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type Bukk struct {
	Source string
}

func New() *Bukk {
	return &Bukk{}
}

func (b *Bukk) SetSource(source string) {
	b.Source = source
}

func (b *Bukk) Name() string {
	return "bukkake.to"
}

var (
	regexSubmatchURL    = regexp.MustCompile(`data-preview="([^"]*)"`)
	regexSubmatchSource = regexp.MustCompile(`<a itemprop="target"[^>]*href="([^"]*)"`)
)

func (b *Bukk) Get(index int) (list []preview.ContentResource, err error) {
	page := index + 1
	uri := `https://bukkake.to/latest-updates/?mode=async&function=get_block&block_id=list_videos_latest_videos_list&sort_by=post_date&from=` + strconv.Itoa(page)
	res, err := resty.New().R().Get(uri)
	submatchURL := regexSubmatchURL.FindAllStringSubmatch(res.String(), -1)
	submatchSource := regexSubmatchSource.FindAllStringSubmatch(res.String(), -1)
	if len(submatchURL) != len(submatchSource) {
		err = fmt.Errorf("len(submatchURL) != len(submatchSource)")
		return
	}
	for index, match := range submatchURL {
		matchSource := submatchSource[index]
		if len(matchSource) < 2 || len(match) < 2 {
			continue
		}
		list = append(list, preview.ContentResource{Source: matchSource[1], View: submatchURL[index][1]})
	}
	return
}

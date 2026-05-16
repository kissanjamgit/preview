// Package aziani something something
package aziani

import (
	"encoding/json"
	"fmt"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type Aziani struct {
	source string
}

func (a Aziani) Name() string {
	return "aziani"
}

func (a *Aziani) SetSource(source string) {
	a.source = source
}

func New() *Aziani {
	return &Aziani{}
}

var header = map[string]string{
	`Host`:                 `azianistudios.com`,
	`User-Agent`:           `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:150.0) Gecko/20100101 Firefox/150.0`,
	`Accept`:               `application/json, text/plain, */*`,
	`Accept-Language`:      `en-US,en;q=0.9`,
	`Accept-Encoding`:      `gzip, deflate, br, zstd`,
	`X-NATS-cms-area-id`:   `3b4c609c-6a0d-4cb9-9cce-0605f32b79ec`,
	`x-nats-natscode`:      `MC4wLjIuMi4wLjAuMC4wLjA`,
	`x-nats-entity-decode`: `1`,
	`Origin`:               `https://aziani.com`,
	`Sec-GPC`:              `1`,
	`Connection`:           `keep-alive`,
	`Referer`:              `https://aziani.com/`,
	`Sec-Fetch-Dest`:       `empty`,
	`Sec-Fetch-Mode`:       `cors`,
	`Sec-Fetch-Site`:       `cross-site`,
	`TE`:                   `trailers`,
}

const size = 18

type jsonAdapte struct {
	Sets []struct {
		CMSID      string `json:"cms_set_id"`
		Name       string `json:"name"`
		Slug       string `json:"slug"`
		PreviewFmt struct {
			Clip []struct {
				Signature string `json:"signature"`
				Fileuri   string `json:"fileuri"`
			} `json:"clip"`
			Trailer struct {
				Formats []struct {
					Content []struct {
						Signature string `json:"signature"`
						Fileuri   string `json:"fileuri"`
					} `json:"content"`
				} `json:"formats"`
			} `json:"trailer"`
		} `json:"preview_formatted"`
	} `json:"sets"`
}

const (
	baseURL = `https://azianistudios.com`
	cmsURL  = `https://c75c0c3063.mjedge.net`
)

var (
	sourceFmt = `https://aziani.com/video/%s`
	viewFmt   = `%s%s?%s`
)

func (a Aziani) Get(index int) (list []preview.ContentResource, err error) {
	uri := fmt.Sprintf(`%s/tour_api.php/content/sets?cms_set_ids=&data_types=1&content_count=1&count=18&start=%d&cms_area_id=3b4c609c-6a0d-4cb9-9cce-0605f32b79ec&cms_block_id=115117&orderby=published_desc&content_video_orientation=horizontal,vertical,square&content_type=video&status=enabled&text_search=&data_type_search=%%7B%%227%%22:%%22436%%22%%7D`, baseURL, size*index)
	client := resty.New()
	res, err := client.R().SetHeaders(header).Get(uri)
	if err != nil {
		return
	}
	var ja jsonAdapte
	json.Unmarshal(res.Bytes(), &ja)
	for _, item := range ja.Sets {
		source := fmt.Sprintf(sourceFmt, item.CMSID)
		// formats := item.PreviewFmt.Trailer.Formats
		clips := item.PreviewFmt.Clip
		if len(clips) < 1 {
			continue
		}
		clip := clips[0]
		// if len(formats) < 1 {
		// 	continue
		// }
		// content := formats[0].Content
		// if len(content) < 1 {
		// 	continue
		// }
		// view := fmt.Sprintf(viewFmt, cmsURL, content[0].Fileuri, content[0].Signature)
		view := fmt.Sprintf(viewFmt, cmsURL, clip.Fileuri, clip.Signature)
		list = append(list, preview.ContentResource{Source: source, View: view})
	}
	return
}

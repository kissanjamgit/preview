// Package teamskt
package teamskt

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/kissanjamgit/preview"

	"resty.dev/v3"
)

type teamskt struct {
	Source string
}

type adaptJSONHits struct {
	Hits struct {
		Hits []struct {
			Source struct {
				ID      string `json:"id"`
				Video   string `json:"video"`
				Title   string `json:"title"`
				Trailer string `json:"videoTrailer"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

var size = 30

var tree = map[string]string{"teamskeet": "ts_network", "swappz": "swap_bundle", "freeuse": "freeusebundle"}

func (t *teamskt) Get(index int) (list []preview.ContentResource, err error) {
	if index < 0 {
		err = fmt.Errorf("index must be greater than 0")
		return
	}
	studioPath := tree[t.Source]
	pad := size * index
	client := resty.New()
	uri := "https://tours-store.psmcdn.net/" + studioPath + "/_search?sort=publishedDate:desc&q=(type:video%20AND%20isXSeries:false%20AND%20isUpcoming:false)&size=" + strconv.Itoa(size) + "&from=" + strconv.Itoa(pad)
	res, err := client.R().Get(uri)
	var adapt adaptJSONHits
	json.Unmarshal(res.Bytes(), &adapt)
	for _, item := range adapt.Hits.Hits {
		// list = append(list, preview.ContentResource{Source: fmt.Sprintf("https://%s.com/movies/%s", t.Source, item.Source.ID), View: item.Source.Trailer})
		list = append(list, preview.ContentResource{Source: item.Source.Video + `@` + strings.Replace(item.Source.Title, " ", "-", -1), View: item.Source.Trailer})
	}
	return
}

var Domain = []string{
	`teamskeet`, `swappz`, `freeuse`,
}

func New(source string) preview.Preview {
	return &teamskt{source}
}

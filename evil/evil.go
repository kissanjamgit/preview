// Package evil do something
package evil

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"regexp"
	"strings"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type evil struct {
	Source string
}

func (e *evil) SetSource(source string) {
	e.Source = source
}

func (*evil) Name() string {
	return Domain[0]
}

func New() *evil {
	return &evil{}
}

var Domain = []string{`evilangel`, `moderndaysins`, `devilsfilm`, `puretaboo`, `mommysboy`, `mommysgirl`, `girlsway`, `21sextury`, `outofthefamily`, `tabooheat`, `accidentalgangbang`, `mommyblowsbest`, `nurumassage`, `filthykings`, `dogfartnetwork`, `gangbangcreampie`, `filthykings`}

type devilsfilmAPIEnv struct {
	API struct {
		Algolia struct {
			AppID  string `json:"applicationID"`
			APIKey string `json:"apiKey"`
		} `json:"algolia"`
	} `json:"api"`
}
type adaptJSONHits struct {
	Results []struct {
		Hits []struct {
			ClipID   int    `json:"clip_id"`
			Sitename string `json:"sitename"`
			URLPath  string `json:"url_title"`
		} `json:"hits"`
	} `json:"results"`
}

type StudioPath struct {
	name string
	path string
}

var Header = map[string]string{
	`User-Agent`:                `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) Gecko/20100101 Firefox/135.0`,
	`Accept`:                    `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`,
	`Accept-Language`:           `en-US,en;q=0.5`,
	`Accept-Encoding`:           `gzip, deflate`,
	`DNT`:                       `1`,
	`Sec-GPC`:                   `1`,
	`Connection`:                `keep-alive`,
	`Upgrade-Insecure-Requests`: `1`,
	`Sec-Fetch-Dest`:            `document`,
	`Sec-Fetch-Mode`:            `navigate`,
	`Sec-Fetch-Site`:            `cross-site`,
	`Priority`:                  `u=0, i`,
	`TE`:                        `trailers`,
}

func body(host string, index int) string {
	if host == `dogfartnetwork` {
		return fmt.Sprintf(`{"requests":[{"indexName":"all_scenes_latest_desc","analytics":true,"analyticsTags":["component:searchlisting","section:freetour","site:%s","context:videos","device:desktop"],"clickAnalytics":true,"facetingAfterDistinct":true,"facets":["categories.url_name","sitename"],"filters":"(NOT categories.name:'Behind The Scene') AND (NOT content_tags:'gay') AND (NOT sitename:'dfxtrapartners') AND (upcoming:'0')","highlightPostTag":"__/ais-highlight__","highlightPreTag":"__ais-highlight__","hitsPerPage":60,"maxValuesPerFacet":1000,"page":%d,"query":""}]}`, host, index)
	}
	return fmt.Sprintf(`{"requests":[{"indexName":"all_scenes_latest_desc","analytics":true,"analyticsTags":["component:searchlisting","section:freetour","site:%s","context:videos","device:desktop"],"clickAnalytics":true,"facetingAfterDistinct":true,"facets":["categories.url_name"],"filters":"(upcoming:'0') AND availableOnSite:%s","highlightPostTag":"__/ais-highlight__","highlightPreTag":"__ais-highlight__","hitsPerPage":60,"maxValuesPerFacet":1000,"page":%d,"query":""}]}`, host, host, index)
}

func (e evil) Get(index int) (list []preview.ContentResource, err error) {
	origin := fmt.Sprintf("https://www.%s.com", e.Source)
	client := resty.New()

	res, err := client.SetHeaders(Header).R().Get(origin)
	if err != nil {
		return
	}
	var match []byte
	if submatch := regexp.MustCompile(`window.env\s+=\s+([^;]+)`).FindSubmatch(res.Bytes()); len(submatch) < 2 {
		err = fmt.Errorf("len(match) < 2, len(match): %d", len(submatch))
		return
	} else {
		match = submatch[1]
	}
	var apiEnv devilsfilmAPIEnv
	err = json.Unmarshal(match, &apiEnv)
	if err != nil {
		return
	}
	u := url.URL{
		Scheme:  "https",
		Host:    strings.ToLower(apiEnv.API.Algolia.AppID) + "-dsn.algolia.net",
		Path:    "/1/indexes/*/queries",
		RawPath: "/1/indexes/*/queries",
	}

	q := u.Query()
	q.Set(`x-algolia-agent`, `Algolia for JavaScript (5.50.1); Lite (5.50.1); Browser; instantsearch.js (4.74.0); react (18.2.0); react-instantsearch (7.13.0); react-instantsearch-core (7.13.0); JS Helper (3.22.4)`)
	q.Set(`x-algolia-api-key`, apiEnv.API.Algolia.APIKey)
	q.Set(`x-algolia-application-id`, apiEnv.API.Algolia.AppID)

	u.RawQuery = q.Encode()

	header := maps.Clone(Header)
	header[`Origin`] = origin
	header[`Referer`] = origin
	res, err = client.R().SetBody(body(e.Source, index)).SetHeaders(header).Post(u.String())
	if err != nil {
		return
	}
	var aj adaptJSONHits
	err = json.Unmarshal(res.Bytes(), &aj)
	if err != nil {
		return
	}

	if len(aj.Results) < 1 {
		err = fmt.Errorf(`len(aj.Results) < 1 `)
		return
	}
	for _, item := range aj.Results[0].Hits {
		source := fmt.Sprintf("%s/en/video/%s/%s/%d", origin, item.Sitename, item.URLPath, item.ClipID)
		view := fmt.Sprintf("https://videothumb.gammacdn.com/500x281/%d.mp4", item.ClipID)
		list = append(list, preview.ContentResource{Source: source, View: view})
	}
	return
}

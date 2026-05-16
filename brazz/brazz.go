// Package brazz something something
package brazz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"regexp"
	"strings"
	"time"

	"github.com/kissanjamgit/preview"
	"resty.dev/v3"
)

type Brazz struct {
	Source string
}

func New() *Brazz {
	return &Brazz{}
}

func (b Brazz) SearchIdentity()       {}
func (b Brazz) SearchDomainIdentity() {}

var domain = []string{`brazzers`, `milfed`, `bangbros`, `realitykings`, `twistys`, `digitalplayground`, `mofos`, `letsdoeit`, `mypervyfamily`, `deviante`}

func (b Brazz) Domain() []string { return domain }

func (b Brazz) Name() string {
	return domain[0]
}

func (b *Brazz) SetSource(source string) {
	b.Source = source
}

func (b Brazz) Get(index int) ([]preview.ContentResource, error) {
	return b.get(nil, index)
}

func (b Brazz) get(query []string, index int) (list []preview.ContentResource, err error) {
	client := resty.New()
	defer client.Close()
	var jwt string

	domain := fmt.Sprintf("https://%s.com/", b.Source)
	req, err := client.R().SetHeaders(header(nil, domain)).Get(domain)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := regexp.MustCompile(`"jwt":"(\S*?)"`).FindStringSubmatch(req.String())
	jwt = data[len(data)-1]

	var url string
	Req := client.R().SetHeaders(header(&jwt, domain))
	if len(query) != 0 {
		url = fmt.Sprintf("https://site-api.project1service.com/v1/dd/videos?pageType=SEARCH_VIDEOS&limit=24&offset=%d&orderBy=newest&query=%s&sexualOrientation=straight&source=p1", 24*index, strings.Join(query, `+`))
		fmt.Println(url)
		res, e := Req.Get(url)
		if err != nil {
			return nil, e
		}
		var ja JSONAdapteQuery
		e = json.Unmarshal(res.Bytes(), &ja)
		if e != nil {
			return nil, e
		}
		fmt.Println(ja)
		for _, item := range ja.Results {
			list = append(list, preview.ContentResource{Source: cleanBrazz(b.Source, item.Title, item.ID), View: item.Videos.Mediabook.Files.Res720p.URL.View})
		}
		return list, nil
	}

	url = fmt.Sprintf("https://site-api.project1service.com/v2/releases?adaptiveStreamingOnly=false&dateReleased=%%3C%s&orderBy=-dateReleased&type=scene&limit=40&offset=%d", time.Now().Format("2006-01-02"), 40*index)
	res, e := Req.Get(url)
	if err != nil {
		return nil, e
	}
	var adapt AdaptarJSONParse
	e = json.Unmarshal(res.Bytes(), &adapt)
	if e != nil {
		return nil, e
	}

	for _, item := range adapt.Result {
		list = append(list, preview.ContentResource{Source: cleanBrazz(b.Source, item.Title, item.ID), View: item.Video.Mediabook.Files.Res720p.Urls.View})
	}
	return
}

func (b Brazz) Search(query []string, index int) ([]preview.ContentResource, error) {
	return b.get(query, index)
}

var sceneOrVideo = map[string]bool{
	"milfed":            true,
	"mofos":             true,
	"twistys":           true,
	"realitykings":      true,
	"digitalplayground": true,
	`letsdoeit`:         true,
	`mypervyfamily`:     false,
}

var (
	invalidChars = regexp.MustCompile(`[^a-z0-9\s]`)
	header_      = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0",
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "en-US,en;q=0.5",
		"Origin":          "https://www.brazzers.com",
		"Sec-GPC":         "1",
		"Connection":      "keep-alive",
		"Referer":         "https://www.brazzers.com/",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "cross-site",
		"TE":              "trailers",
	}
)

func header(jwt *string, domain string) map[string]string {
	h := map[string]string{}
	maps.Copy(h, header_)
	if jwt != nil && *jwt != "" {
		h["instance"] = *jwt
	}
	h["Origin"] = domain
	h["Referer"] = domain
	return h
}

func cleanBrazz(domain string, title string, id int) string {
	s := title
	s = strings.ToLower(s)
	s = invalidChars.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, " ", "-")

	S := "scene"
	if b := sceneOrVideo[domain]; !b {
		S = "video"
	}
	return fmt.Sprintf("https://%s.com/%s/%v/%s", domain, S, id, s)
}

type JSONAdapteQuery struct {
	Results []struct {
		ID     int         `json:"id"`
		Title  string      `json:"title"`
		Videos VideosField `json:"videos"` // Uses the custom type
	} `json:"result"`
}

// VideosField represents the structure when videos are present
type VideosField struct {
	Mediabook struct {
		Files struct {
			Res720p struct {
				URL struct {
					View string `json:"view"`
				} `json:"urls"`
			} `json:"720p"`
		} `json:"files"`
	} `json:"mediabook"`
}

// UnmarshalJSON handles both {} (object) and [] (empty array) variants gracefully
func (v *VideosField) UnmarshalJSON(b []byte) error {
	trimmed := bytes.TrimSpace(b)
	if len(trimmed) == 0 {
		return nil
	}

	// If the JSON starts with '[', it's an empty array "[]".
	// We skip unmarshaling so it retains its default zero values.
	if trimmed[0] == '[' {
		return nil
	}

	// If it's an object '{', unmarshal normally.
	// We use an Alias type to avoid an infinite recursion loop during unmarshaling.
	type Alias VideosField
	var a Alias
	if err := json.Unmarshal(trimmed, &a); err != nil {
		return err
	}

	*v = VideosField(a)
	return nil
}

type AdapatJSONVideo struct {
	Mediabook struct {
		Files struct {
			Res720p struct {
				Urls struct {
					View string `json:"view"`
				} `json:"urls"`
			} `json:"720p"`
		} `json:"files"`
	} `json:"mediabook"`
}

type AdapatJSONResult struct {
	ID    int             `json:"id"`
	Title string          `json:"title"`
	Video AdapatJSONVideo `json:"videos"`
}

type AdaptarJSONParse struct {
	Result []AdapatJSONResult `json:"result"`
}

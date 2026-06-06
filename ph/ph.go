// Package ph is something
package ph

import (
	"fmt"
	"math/rand/v2"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/kissanjamgit/preview"
	"github.com/kissanjamgit/preview/common"
	"golang.org/x/sync/errgroup"
	"resty.dev/v3"
)

type ph struct {
	Source string
}

func New() *ph {
	return &ph{
		Source: ``,
	}
}

func (p *ph) SetSource(source string) {
	p.Source = source
}

func (p *ph) Name() string {
	return "pornhub"
}

const (
	uriShorties = `https://www.pornhub.org/shorties`
	baseURL     = "https://www.pornhub.com"
)

func (*ph) DomainIdentity() {}
func (*ph) Domain() []string {
	return []string{`normal`, `short`}
}

func (p *ph) get(query []string, index int) (list []preview.ContentResource, err error) {
	page := index + 1

	var uri string
	if len(query) == 0 {
		uri = fmt.Sprintf(`%s/video?o=ht&page=%d`, baseURL, page)
	} else {
		uri = fmt.Sprintf(`%s/video/search?search=%s&page=%d`, baseURL, strings.Join(query, `+`), page)
	}
	client := resty.New()
	defer client.Close()
	res, err := client.R().Get(uri)
	if err != nil {
		return
	}
	submatchID := regexp.MustCompile(`data-video-vkey="([^"]*)"`).FindAllStringSubmatch(res.String(), -1)
	submatchURL := regexp.MustCompile(`data-mediabook="([^"]*)"[^>]*title="([^"]*)"`).FindAllStringSubmatch(res.String(), -1)
	if len(submatchID) != len(submatchURL) {
		err = fmt.Errorf(`len(submatchID) != len(submatchURL)`)
		return
	}
	for index, match := range submatchURL[4:] {
		if len(match) < 3 {
			continue
		}
		list = append(list, preview.ContentResource{Source: `https://www.pornhub.com/view_video.php?viewkey=` + submatchID[index][1], View: match[1]})
	}
	return
}

func (p *ph) Get(index int) ([]preview.ContentResource, error) {
	switch p.Source {
	case `normal`:
		return p.get(nil, index)
	case `short`:
		common.PlayerArgs = append(common.PlayerArgs, `--http-header-fields=Referer: https://www.pornhub.org/`, `--ytdl-raw-options=impersonate=chrome`)
		return p.shorties()

	default:
		return nil, fmt.Errorf(`invalid source`)
	}
}

func shortiesGoRoutine(req *resty.Request) (list []preview.ContentResource, err error) {
	res, err := req.Get(uriShorties)
	if err != nil {
		return
	}
	submatchShortie := regexp.MustCompile(`"linkUrl":"([^"]*)","shortieUrl":"([^"]*)"`).FindAllStringSubmatch(res.String(), -1)
	submatchVideoURL := regexp.MustCompile(`"videoUrl":"([^"]*)",\s*"quality":"480"`).FindAllStringSubmatch(res.String(), -1)

	if len(submatchShortie) < 1 || len(submatchVideoURL) < 1 {
		err = fmt.Errorf(`len(submatchShortie) < 1 || len(submatchVideoURL) < 1`)
		return
	}

	for matchIndex, match := range submatchShortie {
		if len(match) < 3 || len(submatchVideoURL[matchIndex]) < 2 {
			continue
		}
		source, err := url.Parse(regexpEscapeForwardSlash.ReplaceAllString(match[1], `/`))
		if err != nil {
			continue
		}
		view, err := url.Parse(regexpEscapeForwardSlash.ReplaceAllString(submatchVideoURL[matchIndex][1], `/`))
		if err != nil {
			continue
		}
		list = append(list, preview.ContentResource{Source: source.String(), View: view.String()})
	}

	return
}

var regexpEscapeForwardSlash = regexp.MustCompile(`\\/`)

func extTokenAndSeed(s string) (string, string) {
	seed := regexp.MustCompile(`ddseed = \d+`).FindString(s)
	return regexp.MustCompile(`var token = "([^"]*)"`).FindString(s), seed
}

func (p *ph) shorties() (list []preview.ContentResource, err error) {
	// goRoutine doen't work as expected due to too little time delay all the goRoutine recives the same data, os better live it with sync
	var g errgroup.Group
	goRoutineNo := 5
	goRoutineList := make([][]preview.ContentResource, goRoutineNo)

	client := resty.New()
	defer client.Close()
	res, err := client.R().Get(uriShorties)
	if err != nil {
		return
	}
	token, seed := extTokenAndSeed(res.String())
	// https://www.pornhub.org/shorties
	// uri := url.URL{Host: `www.pornhub.org`, Path: `shorties`, Scheme: "https"}
	for index := range goRoutineNo {
		i := index

		q := url.Values{}
		q.Set(`token`, token)
		q.Set(`seed`, seed)
		q.Set(`orientation`, `straight`)
		q.Set(`offset`, strconv.Itoa(rand.IntN(5000)))
		g.Go(
			func() (er error) {
				req := client.R().SetQueryParamsFromValues(q)

				resList, er := shortiesGoRoutine(req)
				if err != nil {
					return er
				}
				goRoutineList[i] = resList
				return nil
			},
		)
	}

	err = g.Wait()

	for _, item := range goRoutineList {
		list = append(list, item...)
	}
	return
}

func (p *ph) SearchIdentity() {}

func (p *ph) Search(query []string, index int) ([]preview.ContentResource, error) {
	return p.get(query, index)
}

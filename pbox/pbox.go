// Package pbox is a preview provider for
package pbox

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kissanjamgit/preview"

	"github.com/tidwall/gjson"
	"resty.dev/v3"
)

type PBox struct {
	source string
}

type contentResult struct {
	contentID      string
	runtime        time.Duration
	videoPreview   string
	altercontentID string
}

func getContentResult(_, value gjson.Result, result map[string]contentResult, sortResult *[]string) bool {
	id := value.Get("id").String()
	var runtime time.Duration
	runtimeString := value.Get("runtime").String()
	if runtimeString == "" {
		return true
	} else {
		runtime = parseDuration(runtimeString)
	}
	videoPreview := value.Get("video_preview").String()
	alterContentID := value.Get("alternate.content_id").String()
	result[id] = contentResult{contentID: id, runtime: runtime, videoPreview: videoPreview, altercontentID: alterContentID}
	*sortResult = append(*sortResult, id)
	return true
}

func parseDuration(str string) time.Duration {
	parts := strings.Split(str, ":")
	if len(parts) != 3 {
		return 0
	}

	// Convert strings to integers
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	sVal, _ := strconv.Atoi(parts[2])

	// Calculate total duration
	// h*Hour + m*Minute + s*Second
	return time.Duration(h)*time.Hour +
		time.Duration(m)*time.Minute +
		time.Duration(sVal)*time.Second
}

func removeAlterContentWithSmallDuration(dict map[string]contentResult) {
	toDelete := make([]string, 0)
	for k, v := range dict {
		altercontentID := v.altercontentID
		if altercontentID == "" {
			continue
		}
		if dict[altercontentID].runtime < v.runtime {
			toDelete = append(toDelete, altercontentID)
			continue
		}
		toDelete = append(toDelete, k)
	}
	for _, k := range toDelete {
		delete(dict, k)
	}
}

func handleStudio(client *resty.Client, id int, index int) (list []preview.ContentResource, err error) {
	res, err := client.R().Get(fmt.Sprintf("https://pornbox.com/studio/%d/?skip=%d&sort=latest", id, index))
	if err != nil {
		fmt.Println(err)
		return
	}
	result := make(map[string]contentResult)
	var sortResult []string
	contentsList := gjson.Get(res.String(), "contents")
	contentsList.ForEach(func(key, value gjson.Result) bool {
		return getContentResult(key, value, result, &sortResult)
	})
	removeAlterContentWithSmallDuration(result)

	for _, id := range sortResult {
		v, ok := result[id]
		if !ok {
			continue
		}
		list = append(list, preview.ContentResource{Source: "https://pornbox.com/application/watch-page/" + v.contentID, View: v.videoPreview})

	}
	return
}

var Domain = func() (l []string) {
	for k := range domain {
		l = append(l, k)
	}
	sort.Strings(l)
	return
}()

var domain = map[string]int{
	"AdventureTeens":      1007,
	"GiorgioGrandi":       29,
	"GiorgioLabs":         33,
	"LancelotStyles":      698,
	"Angelogodshack":      251,
	"YummyeStudio":        111,
	"Disciples Of Desire": 1214,
	"NatashateenFilms":    74,
	"AnalTerror":          23041,
	"NRX-Studio":          94,
	"Mambo Perv":          197,
}

func New(source string) preview.Preview {
	return &PBox{source: source}
}

func (p *PBox) Get(index int) (list []preview.ContentResource, err error) {
	client := resty.New()
	value, ok := domain[p.source]
	if !ok {
		err = fmt.Errorf("invalid domain")
		return
	}

	list, err = handleStudio(client, value, index)
	return
}

package faphouse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FaphousePair struct {
	View  string `json:"view"`
	Source string `json:"source"`
}

type FaphouseResponse struct {
	Pairs []FaphousePair `json:"pairs"`
}

type Faphouse struct {
	URL string
	Pairs []FaphousePair
	StudioPath
}

type StudioPath struct {
	name string
	path string
}

func (f *Faphouse) Fetch() ([]FaphousePair, error) {
	resp, err := http.Get(f.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var data FaphouseResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return parseHTMLPairs(string(body))
	}

	f.Pairs = data.Pairs
	return f.Pairs, nil
}

func (f *Faphouse) FetchHTML() ([]FaphousePair, error) {
	resp, err := http.Get(f.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	f.Pairs, err = parseHTMLPairs(string(body))
	return f.Pairs, err
}

func (f *Faphouse) FetchFile() ([]FaphousePair, error) {
	data, err := os.ReadFile(f.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	pairs, err := parseHTMLPairs(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %v", err)
	}

	f.Pairs = pairs
	return f.Pairs, nil
}

func parseHTMLPairs(htmlContent string) ([]FaphousePair, error) {
	var pairs []FaphousePair

	// Regex pattern to match view and source attributes in HTML
	viewRegex := regexp.MustCompile(`(?i)view\s*=\s*["']([^"']+)["']`)
	sourceRegex := regexp.MustCompile(`(?i)source\s*=\s*["']([^"']+)["']`)

	lines := strings.Split(htmlContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "view") && strings.Contains(line, "source") {
			viewMatch := viewRegex.FindStringSubmatch(line)
			sourceMatch := sourceRegex.FindStringSubmatch(line)

			var view, source string
			if len(viewMatch) > 1 {
				view = viewMatch[1]
			}
			if len(sourceMatch) > 1 {
				source = sourceMatch[1]
			}

			if view != "" && source != "" {
				pair := FaphousePair{
					View:  view,
					Source: source,
				}
				pairs = append(pairs, pair)
			}
		}
	}

	return pairs, nil
}

func (f *Faphouse) GetView() string {
	return f.URL
}

func (f *Faphouse) GetSource() string {
	if len(f.Pairs) > 0 {
		return f.Pairs[0].Source
	}
	return "faphouse"
}

func (f *Faphouse) GetName() string {
	return f.name
}

func (f *Faphouse) GetPath() string {
	return f.path
}

func NewFaphouse(url string, name string) *Faphouse {
	return &Faphouse{
		URL: url,
		Pairs: []FaphousePair{},
		StudioPath: StudioPath{
			name: name,
			path: "faphouse",
		},
	}
}

func main() {
	f := NewFaphouse("http://example.com/faphouse.html", "faphouse")
	pairs, err := f.FetchHTML()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, pair := range pairs {
		fmt.Printf("View: %s, Source: %s\n", pair.View, pair.Source)
	}

	fmt.Printf("Total pairs found: %d\n", len(pairs))
}

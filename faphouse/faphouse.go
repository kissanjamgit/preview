package faphouse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func parseHTMLPairs(htmlContent string) ([]FaphousePair, error) {
	var pairs []FaphousePair

	// Simple regex-like parsing for view/source pairs in HTML
	lines := strings.Split(htmlContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "view") && strings.Contains(line, "source") {
			pair := FaphousePair{
				View:  extractValue(line, "view"),
				Source: extractValue(line, "source"),
			}
			if pair.View != "" && pair.Source != "" {
				pairs = append(pairs, pair)
			}
		}
	}

	return pairs, nil
}

func extractValue(line, key string) string {
	parts := strings.SplitN(line, key+":", 2)
	if len(parts) < 2 {
		return ""
	}
	value := strings.TrimSpace(parts[1])
	if idx := strings.Index(value, ","); idx != -1 {
		value = value[:idx]
	}
	return strings.TrimSpace(value)
}

func (f *Faphouse) GetView() string {
	return f.URL
}

func (f *Faphouse) GetSource() string {
	// Default source can be configured or extracted from URL
	return "faphouse"
}

func NewFaphouse(url string) *Faphouse {
	return &Faphouse{URL: url}
}

func main() {
	f := NewFaphouse("http://example.com/faphouse.html")
	pairs, err := f.FetchHTML()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, pair := range pairs {
		fmt.Printf("View: %s, Source: %s\n", pair.View, pair.Source)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

// Resource an MDN page
type Resource struct {
	ID           uint
	Label        string
	Locale       string
	Modified     string
	Slug         string
	Subpages     []Resource
	Summary      string
	Tags         []string
	Title        string
	Translations []Resource
	UUID         string
	URL          string
}

type SearchSettings struct {
	RankingRules         []string            `json:"rankingRules,omitempty"`
	RankingDistinct      string              `json:"rankingDistinct,omitempty"`
	SearchableAttributes []string            `json:"searchableAttributes,omitempty"`
	DisplayedAttributes  []string            `json:"displayedAttributes,omitempty"`
	StopWords            []string            `json:"stopWords,omitempty"`
	Synonyms             map[string][]string `json:"synonyms,omitempty"`
	IndexNewFields       bool                `json:"indexNewFields,omitempty"`
}

type SearchResponse struct {
	Hits             []Resource `json:"hits"`
	Offset           uint       `json:"offset"`
	Limit            uint       `json:"limit"`
	ProcessingTimeMs uint       `json:"processingTimeMs"`
	Query            string
}

type Loader struct {
	WG             *sync.WaitGroup
	SearchIndex    string
	SearchEndpoint string
	Errors         chan error
}

var ErrUnexpectedResponse = errors.New("Unexpected response from API")

func (l *Loader) loadResource(resource Resource) {
	l.WG.Add(1)
	defer l.WG.Done()

	body, err := json.Marshal([]Resource{resource})
	if err != nil {
		l.Errors <- err
		return
	}

	res, err := http.Post(l.indexURL()+"/documents", "application/json", bytes.NewReader(body))
	if err != nil {
		l.Errors <- err
		return
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		l.Errors <- ErrUnexpectedResponse
	}

	for _, sub := range resource.Subpages {
		go l.loadResource(sub)
	}
}

func (l *Loader) loadEndpoint(endpoint string) {
	res, err := http.Get("https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/" + endpoint + "$children?expand")
	if err != nil {
		l.Errors <- err
		return
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		l.Errors <- ErrUnexpectedResponse
		return
	}

	decoder := json.NewDecoder(res.Body)
	for decoder.More() {
		var p Resource
		err = decoder.Decode(&p)
		if err != nil {
			l.Errors <- err
			return
		}

		go l.loadResource(p)
	}
}

func (l *Loader) load() {
	endpoints := []string{
		"Global_Objects",
		"Operators",
		"Statements",
		"Functions",
		"Classes",
		"Errors",
	}

	for _, e := range endpoints {
		go l.loadEndpoint(e)
	}

	return
}

func (l *Loader) setup() (err error) {
	data, err := json.Marshal(struct {
		UID string `json:"uid"`
	}{l.SearchIndex})
	if err != nil {
		return
	}

	_, err = http.Post(l.SearchEndpoint+"/indexes", "application/json", bytes.NewReader(data))
	if err != nil {
		return
	}
	// ignore response code

	data, err = json.Marshal(&SearchSettings{
		Synonyms: map[string][]string{
			".prototype.": []string{"#"},
			"#":           []string{".prototype."},
		},
		SearchableAttributes: []string{
			"Title",
			"Tags",
			"Label",
			"Summary",
		},
	})
	if err != nil {
		return
	}

	res, err := http.Post(l.indexURL()+"/settings", "application/json", bytes.NewReader(data))
	if err != nil {
		return
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = ErrUnexpectedResponse
	}
	return
}

func (l *Loader) indexURL() string {
	return l.SearchEndpoint + "/indexes/" + l.SearchIndex
}

func (l *Loader) search(query string) (out Resource, err error) {
	params := url.Values{
		"q":     []string{query},
		"limit": []string{strconv.Itoa(1)},
	}

	res, err := http.Get(l.indexURL() + "/search?" + params.Encode())
	if err != nil {
		return
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = ErrUnexpectedResponse
		return
	}

	results := &SearchResponse{}
	err = json.NewDecoder(res.Body).Decode(results)
	if len(results.Hits) > 0 {
		out = results.Hits[0]
	}
	return
}

package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/appellation/mdn/sonic"
)

const (
	// CollectionSummary is the collection key of the MDN summaries
	CollectionSummary = "js_summary"

	// Bucket is the Sonic bucket things are stored in
	Bucket = "default"
)

// Resource an MDN page
type Resource struct {
	ID           uint64
	Label        string
	Locale       string
	Modified     string
	Slug         string
	Subpages     []*Resource
	Summary      string
	Tags         []string
	Title        string
	Translations []*Resource
	UUID         string
	URL          string
}

var endpoints = []string{
	"Global_Objects",
	"Operators",
	"Statements",
	"Functions",
	"Classes",
	"Errors",
}

var resources = sync.Map{}

func ingestResource(ingest *sonic.Ingest, resource *Resource) error {
	id := strconv.FormatUint(resource.ID, 10)
	resources.Store(id, resource)

	if len(resource.Summary) > 0 {
		res, err := http.Get("https://developer.mozilla.org" + resource.URL)
		if err != nil {
			return err
		}

		switch res.StatusCode {
		case http.StatusOK:
		default:
			time.Sleep(5 * time.Second)
			return ingestResource(ingest, resource)
		}

		content := &strings.Builder{}
		err = parseContent(res.Body, content)
		if err != nil {
			return err
		}

		res.Body.Close()
		err = ingest.Push(CollectionSummary, Bucket, id, content.String(), "eng")
		if err != nil {
			return err
		}
	}

	wg := sync.WaitGroup{}
	for _, sub := range resource.Subpages {
		wg.Add(1)
		go func(pg *Resource) {
			defer wg.Done()
			err := ingestResource(ingest, pg)
			if err != nil {
				log.Fatal(err)
			}
		}(sub)
	}

	wg.Wait()
	return nil
}

func load(ingester *sonic.Ingest) (err error) {
	ingester.FlushCollection(CollectionSummary)

	wg := sync.WaitGroup{}
	for _, e := range endpoints {
		res, err := http.Get("https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/" + e + "$children?expand")
		if err != nil {
			return err
		}

		p := &Resource{}
		err = json.NewDecoder(res.Body).Decode(p)
		res.Body.Close()
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			ingestResource(ingester, p)
		}()
	}

	wg.Wait()
	return nil
}

func search(search *sonic.Search, query string) (out Resource, err error) {
	results, err := search.Query(CollectionSummary, Bucket, query, sonic.QueryOptions{
		Limit: 1,
	})
	if err != nil {
		return
	}

	log.Println(results)
	if len(results) > 0 {
		ld, _ := resources.Load(results[0])
		out = *ld.(*Resource)
	} else {
		err = errors.New("no results")
	}

	return
}

package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"

	"github.com/arbovm/levenshtein"
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

var keyed = make(map[string]Resource)

func loadResource(base map[string]Resource, resource Resource) {
	base[resource.Title] = resource

	for i := 0; i < len(resource.Subpages); i++ {
		loadResource(base, resource.Subpages[i])
	}
}

func load() (err error) {
	endpoints := []string{
		"Global_Objects",
		"Operators",
		"Statements",
		"Functions",
		"Classes",
		"Errors",
	}

	raw := make(map[string]Resource)
	for _, e := range endpoints {
		res, err := http.Get("https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/" + e + "$children?expand")
		if err != nil {
			return err
		}

		decoder := json.NewDecoder(res.Body)
		for decoder.More() {
			var p Resource
			err = decoder.Decode(&p)
			if err != nil {
				return err
			}

			loadResource(raw, p)
		}
	}

	for key, value := range raw {
		if strings.LastIndex(key, "()") == len(key)-2 {
			key = strings.Replace(key, "()", "", -1)
			keyed[key] = value
		}

		if strings.Contains(key, ".prototype.") {
			key = strings.Replace(key, ".prototype.", "#", -1)
		}

		keyed[key] = value
	}

	return nil
}

func search(query string) (out Resource) {
	var match Resource
	var matchDist = math.MaxInt32

	for title, resource := range keyed {
		distance := levenshtein.Distance(title, query)
		if distance < matchDist {
			match = resource
			matchDist = distance
		}
	}

	return match
}

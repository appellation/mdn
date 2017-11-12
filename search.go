package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/schollz/closestmatch"
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
var titles []string
var resources []Resource

var matcher *closestmatch.ClosestMatch

func loadResource(resource Resource) {
	keyed[resource.Title] = resource

	for i := 0; i < len(resource.Subpages); i++ {
		loadResource(resource.Subpages[i])
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

			loadResource(p)
		}
	}

	length := len(keyed)
	titles = make([]string, length)
	resources = make([]Resource, length)

	i := 0
	j := 0
	for key, value := range keyed {
		if strings.Contains(key, ".prototype.") {
			titles[j] = strings.Replace(key, ".prototype.", "#", -1)
			titles = append(titles, key)
			j += 2
		} else {
			titles[j] = key
			j++
		}

		resources[i] = value
		i++
	}

	matcher = closestmatch.New(titles, []int{2, 3, 4})
	return nil
}

func search(query string) (out Resource) {
	match := matcher.Closest(query)
	return keyed[match]
}

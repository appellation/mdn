package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

func main() {
	errors := make(chan error)
	loader := &Loader{
		WG:             &sync.WaitGroup{},
		SearchIndex:    "js",
		SearchEndpoint: "http://localhost:7700",
		Errors:         errors,
	}
	err := loader.setup()
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})

	loader.load()
	go func() {
		loader.WG.Wait()
		done <- struct{}{}
	}()

	select {
	case err = <-errors:
		panic(err)
	case <-done:
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()["q"]

		if len(q) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "no query")
			return
		}

		query := q[0]
		result, err := loader.search(query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "error parsing response JSON")
			return
		}

		bytes, _ := json.Marshal(result)
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	})

	http.ListenAndServe(":8080", nil)
}

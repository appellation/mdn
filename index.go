package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	load()
	fmt.Println(search("string.prototype").URL)

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()["q"]

		if len(q) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no query"))
			return
		}

		query := q[0]
		bytes, err := json.Marshal(search(query))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "error parsing response JSON")
			return
		}

		w.Write(bytes)
	})

	http.ListenAndServe(":8080", nil)
}

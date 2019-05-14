package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/appellation/mdn/sonic"
)

var (
	sonicHost     = os.Getenv("SONIC_HOST")
	sonicPassword = os.Getenv("SONIC_PASSWORD")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if sonicHost == "" {
		sonicHost = "sonic:1491"
	}
	if sonicPassword == "" {
		sonicPassword = "SecretPassword"
	}
}

func main() {
	conn, err := sonic.Connect(sonicHost, sonicPassword)
	for err != nil {
		time.Sleep(5 * time.Second)
		conn, err = sonic.Connect(sonicHost, sonicPassword)
	}

	ingester, err := conn.Ingest()
	for err != nil {
		log.Fatal(err)
	}

	conn, err = sonic.Connect(sonicHost, sonicPassword)
	if err != nil {
		log.Fatal(err)
	}

	searcher, err := conn.Search()
	if err != nil {
		log.Fatal(err)
	}

	err = load(ingester)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		if len(q) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "no query")
			return
		}

		res, err := search(searcher, q)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}

		log.Println(res)
		bytes, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "error parsing response JSON")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	})

	http.ListenAndServe(":8080", nil)
}

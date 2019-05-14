package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	res, err := http.Get("https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Infinity")
	if err != nil {
		t.Fatal(err)
	}

	content := &strings.Builder{}
	parseContent(res.Body, content)
	res.Body.Close()

	t.Log(content)
}

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "# comment")
		fmt.Fprintln(w, "192.0.2.0")
	}))
	defer ts.Close()
	list, err := Fetch(ts.URL, "http://example.com")
	if err != nil {
		t.Fatal(err)
	}
	got := <-list
	if got.String() != "192.0.2.0" {
		t.Error("got", got)
	}
}

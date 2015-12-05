package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_main(t *testing.T) {
	_, _, _, _, handler := parseFlags()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/v1.0/blocks")
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	// empty array response
	exp := "[]"

	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}

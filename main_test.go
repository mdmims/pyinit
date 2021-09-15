package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var GetList = getList

func TestGetList(t *testing.T) {
	// generate mock response data
	listData := []byte("macos,vscode,python,go,bash")

	// create mock api
	fakeServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "%s", string(listData))
			},
		),
	)
	defer fakeServer.Close()

	testURL := fakeServer.URL

	t.Run(
		"test fetches correct data", func(t *testing.T) {
			got, err := GetList(testURL)
			if err != nil {
				t.Fatalf("Error: %s\n", err)
			}

			if !reflect.DeepEqual(got, listData) {
				t.Errorf("got %q expected %q", got, listData)
			}
		},
	)
}

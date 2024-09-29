package main

import (
	"errors"
	"fmt"
	"github.com/mixailo/go-training-metrics/internal/metrics"
	"github.com/mixailo/go-training-metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

func validatePost(r *http.Request) (status int, err error) {
	if r.Method != http.MethodPost {
		return http.StatusBadRequest, errors.New("invalid method, only POST is allowed")
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		return http.StatusBadRequest, errors.New("only text/plain is supported")
	}

	return
}

func parseURL(url string) (parsed parsedRequestURL, err error) {
	// split URL to get fragments
	fragments := strings.Split(url, "/")

	if len(fragments) != 5 {
		// simple integrity check (TODO: replace by regular expressions or any other suitable method)
		return parsed, errors.New("invalid url, must be 5 fragments")
	}

	if len(fragments[3]) < 1 {
		// empty name
		return parsed, errors.New("invalid item name")
	}

	parsed.itemType = fragments[2]
	parsed.itemName = fragments[3]
	parsed.unConvertedValue = fragments[4]

	return
}

type parsedRequestURL struct {
	itemType         string
	itemName         string
	unConvertedValue string
}

func updateItemValue(w http.ResponseWriter, r *http.Request) {
	// check request method and content-type
	status, err := validatePost(r)
	if err != nil {
		w.WriteHeader(status)
		return
	}

	parsedURL, err := parseURL(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parsedURL.itemType == metrics.TypeCounter.String() {
		// counter type increments stored value
		convertedValue, err := strconv.ParseInt(parsedURL.unConvertedValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		serverStorage.UpdateCounter(parsedURL.itemName, convertedValue)
	} else if parsedURL.itemType == metrics.TypeGauge.String() {
		// gauge type replaces stored value
		convertedValue, err := strconv.ParseFloat(parsedURL.unConvertedValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		serverStorage.UpdateGauge(parsedURL.itemName, convertedValue)
	} else {
		// unknown type
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

var serverStorage = storage.NewStorage()

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", updateItemValue)
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

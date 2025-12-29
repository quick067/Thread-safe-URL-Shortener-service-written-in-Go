package main

import (
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	urlShortener := URLShortener{
		urlStore: make(map[string]string),
	}
	urlShortener.loadFromFile()

	http.HandleFunc("/save", urlShortener.saveURL)
	http.HandleFunc("/", urlShortener.redirectURL)

	http.ListenAndServe(":8080", nil)
}

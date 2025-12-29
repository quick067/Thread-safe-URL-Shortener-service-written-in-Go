package main

import (
	"fmt"
	"io"
	"net/http"
)

func (URLs *URLShortener) saveURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Cannot read body of your request", http.StatusBadRequest)
			return
		}

		longAdress := string(data)
		generatedKey := keyGenerator()

		URLs.mutex.Lock()
		URLs.urlStore[generatedKey] = longAdress
		err = URLs.saveToFile()
		URLs.mutex.Unlock()

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		responseAdress := fmt.Sprintf("http://localhost:8080/%s", generatedKey)
		w.Write([]byte(responseAdress))

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (URLs *URLShortener) redirectURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		path := r.URL.Path
		if len(path) < 1{
			http.Error(w, "Invalid key", http.StatusBadRequest)
		}
		path = path[1:]

		URLs.mutex.RLock()
		value, ok := URLs.urlStore[path]
		URLs.mutex.RUnlock()

		if ok {
			http.Redirect(w, r, value, http.StatusFound)
		} else {
			http.Error(w, "Key didn't found", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

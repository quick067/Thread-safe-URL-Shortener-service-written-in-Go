package handlers

import (
	"URL-shortener/internal/store"
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	Storage *store.Store
}

func (h *Handler) RegisterRoutes() {
	http.HandleFunc("/save", h.saveURL)
	http.HandleFunc("/", h.redirectURL)
}

func (h *Handler) saveURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Cannot read body of your request", http.StatusBadRequest)
			return
		}

		longAdress := string(data)
		generatedKey := keyGenerator()

		if err := h.Storage.SetPair(generatedKey, longAdress); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		responseAdress := fmt.Sprintf("http://localhost:8080/%s", generatedKey)
		w.Write([]byte(responseAdress))

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) redirectURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		path := r.URL.Path[1:]
		if len(path) < 1 {
			http.Error(w, "Invalid key", http.StatusBadRequest)
		}

		value, ok := h.Storage.GetPair(path)

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

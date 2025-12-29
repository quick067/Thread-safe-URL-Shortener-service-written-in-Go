package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const TestGoogleURL = "https://google.com"
const TestLocalURL = "http://localhost:8080"

func TestSaveUrl(t *testing.T) {
	defer os.Remove(fileName)
	s := URLShortener{
		urlStore: make(map[string]string),
	}

	req := httptest.NewRequest("POST", TestLocalURL, strings.NewReader(TestGoogleURL))

	w := httptest.NewRecorder()

	s.saveURL(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected: %d, got: %v", http.StatusOK, w.Code)
	}

	respBody := w.Body.String()
	if !strings.Contains(respBody, TestLocalURL) {
		t.Errorf("Expected: %s, got: %s", TestLocalURL, respBody)
	}

	if len(s.urlStore) == 0 {
		t.Errorf("Expected: >%d elements, got: %d", 0, len(s.urlStore))
	}
}

func TestRedirectUrl(t *testing.T) {
	s := URLShortener{
		urlStore: make(map[string]string),
	}

	s.urlStore["testKey"] = TestGoogleURL

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/testKey", nil)

	s.redirectURL(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected: %v, got: %v", http.StatusFound, w.Code)
	}
	loc := w.Header().Get("Location")
	if loc != TestGoogleURL {
		t.Errorf("Expected: %s, got: %s", TestGoogleURL, loc)
	}
}

package handlers

import (
	"URL-shortener/internal/store"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const TestGoogleURL = "https://google.com"
const TestLocalURL = "http://localhost:8080"

const testFile = "test_storage.json"

var s= store.New(testFile)

var handler Handler = Handler{
	Storage: s,
}

func TestSaveUrl(t *testing.T) {

	defer os.Remove(testFile)

	req := httptest.NewRequest("POST", TestLocalURL, strings.NewReader(TestGoogleURL))

	w := httptest.NewRecorder()

	handler.saveURL(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected: %d, got: %v", http.StatusOK, w.Code)
	}

	respBody := w.Body.String()
	if !strings.Contains(respBody, TestLocalURL) {
		t.Errorf("Expected: %s, got: %s", TestLocalURL, respBody)
	}

	parts := strings.Split(respBody, "/")
	if len(parts) == 0{
		t.Fatal("Response is empty")
	}

	key := parts[len(parts)-1]
	val, ok := s.GetPair(key)
	if !ok{
		t.Error("Key was not saved in storage")
	}else if val != TestGoogleURL{
		t.Errorf("Expected value: %s, got: %s", TestGoogleURL, val)
	}
}

func TestRedirectUrl(t *testing.T) {
	if err := s.SetPair("testKey", TestGoogleURL); err != nil{
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer os.Remove(testFile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/testKey", nil)

	handler.redirectURL(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected: %v, got: %v", http.StatusFound, w.Code)
	}
	loc := w.Header().Get("Location")
	if loc != TestGoogleURL {
		t.Errorf("Expected: %s, got: %s", TestGoogleURL, loc)
	}
}

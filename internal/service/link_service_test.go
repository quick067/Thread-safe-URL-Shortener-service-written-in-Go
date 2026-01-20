package service

import (
	"errors"
	"testing"
	"time"

	"URL-shortener/internal/mocks"
	"URL-shortener/internal/store"

)

func TestCreateLink_Success(t *testing.T){
	ms := &mocks.MockStorage{
		SetPairFunc: func(key, value string, userID int, expiresAt *time.Time) error { return nil},
	}
	ls := &LinkService{Storage: ms, BaseURL: "http://localhost:8080"}

	result, err := ls.CreateLink("https://google.com", "testAlias", "1h", 0)

	if err != nil{
		t.Errorf("Expected: no errors, got: %v", err)
	}
	if result == ""{
		t.Error("Expected: generated link, got: empty string")
	}
}

func TestCreateLink_Conflict(t *testing.T) {
	ms := &mocks.MockStorage{
		SetPairFunc: func(key, value string, userID int, expiresAt *time.Time) error {return store.ErrConflict},
	}
	ls := &LinkService{Storage: ms, BaseURL: "http://localhost:8080"}

	_, err := ls.CreateLink("https://google.com", "testAlias", "1h", 0)

	if err == nil {
		t.Errorf("Expected: %v, got: %v", store.ErrConflict, err)
	}
	if !errors.Is(err, store.ErrConflict){
		t.Errorf("Expected: %v, got: %v", store.ErrConflict, err)
	}
}

func TestCreateUser_InvalidURL(t *testing.T){
	ms := &mocks.MockStorage{
		SetPairFunc: func(key, value string, userID int, expiresAt *time.Time) error { return nil},
	}
	ls := &LinkService{Storage: ms, BaseURL: "http://localhost:8080"}

	_, err := ls.CreateLink("", "testAlias", "1h", 0)

	if err == nil {
		t.Errorf("Expected: %v, got: %v", ErrInvalidURL, err)
	}
	if !errors.Is(err, ErrInvalidURL) {
		t.Errorf("Expected: %v, got: %v", ErrInvalidURL, err)
	}
}

func TestGetURL_Success(t *testing.T) {
	futureTime := time.Now().Add(time.Hour * 999)
	ms := &mocks.MockStorage{
		GetPairFunc: func(path string) (string, *time.Time, error) {return "http://testexample.com", &futureTime, nil},
	}
	ls := &LinkService{Storage: ms, BaseURL: "http://localhost:8080"}

	result, err := ls.GetURL("testAlias")

	if err != nil {
		t.Errorf("Expected: no errors, got: %v", err)
	}
	if result != "http://testexample.com" {
		t.Errorf("Expected: %s, got: %s", "http://testexample.com", result)
	}
}

func TestGetURL_NotFound(t *testing.T) {
	ms := &mocks.MockStorage{
		GetPairFunc: func(path string) (string, *time.Time, error) {return "", nil, store.ErrNotFound},
	}
	ls := &LinkService{Storage: ms, BaseURL: "http://localhost:8080"}

	_, err := ls.GetURL("testAlias")

	if err == nil {
		t.Errorf("Expected: %v, got: %v", store.ErrNotFound, err)
	}
	if !errors.Is(err, store.ErrNotFound){
		t.Errorf("Expected: %v, got: %v", store.ErrNotFound, err)
	}
}

func TestGetURL_Expired(t *testing.T) {
	pastTime := time.Now().Add(-1 *time.Hour)
	ms := &mocks.MockStorage{
		GetPairFunc: func(path string) (string, *time.Time, error) {return "", &pastTime, nil},
	}
	ls := &LinkService{Storage: ms, BaseURL: "http://localhost:8080"}

	_, err := ls.GetURL("testAlias")

	if err == nil {
		t.Errorf("Expected: %v, got: %v", ErrLinkExpired, err)
	}

	if !errors.Is(err, ErrLinkExpired){
		t.Errorf("Expected: %v, got: %v", ErrLinkExpired, err)
	}
}

func TestSaveRandomKey_RetrySuccess(t *testing.T) {
	counter := 0
	futureTime := time.Now().Add(time.Hour * 999)
	ms := &mocks.MockStorage{
		SetPairFunc: func(key, value string, userID int, expiresAt *time.Time) error {
			counter++
			if counter == 1 {
				return store.ErrConflict
			}
			return nil
		},
	}
	ls := &LinkService{Storage: ms, BaseURL: "http://localhost:8080"}

	key, err := ls.saveRandomKey("https://google.com", 0, &futureTime)

	if err != nil {
		t.Errorf("Expected: no errors, got: %v", err)
	}
	if key == "" {
		t.Errorf("Expected: unempty key, got: %s", key)
	}
	if counter != 2 {
		t.Errorf("Expected: counter = %d, got: %d", 2, counter)
	}
}
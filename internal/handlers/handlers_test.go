package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"URL-shortener/internal/mocks"
	"URL-shortener/internal/service"
	"URL-shortener/internal/store"

	"golang.org/x/crypto/bcrypt"
)

func TestSaveURL_Success(t *testing.T) {
	ms := &mocks.MockStorage{
		SetPairFunc: func(key, value string, userID int, expiresAt *time.Time) error {return nil},
	}
	ls := &service.LinkService{
		Storage: ms,
		BaseURL: "http://localhost:8080",
	}
	testHandler := Handler{LinkService: ls}

	reqBody := strings.NewReader(`{"url": "https://google.com", "alias": "test", "duration":"1h"}`)
	req := httptest.NewRequest("POST", "/save", reqBody)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), "user_id", 1.0)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()
	testHandler.saveURL(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected: %v, got: %v", http.StatusOK, recorder.Code)
	}

	if recorder.Body.String() == "" {
		t.Errorf("Expected: non empty string, got: %s", recorder.Body.String())
	}
}

func TestSaveURL_InvalidURL(t *testing.T){
	ms := &mocks.MockStorage{
		SetPairFunc: func(key, value string, userID int, expiresAt *time.Time) error {return service.ErrInvalidURL},
	}
	ls := &service.LinkService{
		Storage: ms,
		BaseURL: "http://localhost:8080",
	}
	testHandler := Handler{LinkService: ls}

	reqBody := strings.NewReader(`{"url": "https://google.com", "alias": "test", "duration":"1h"}`)
	req := httptest.NewRequest("POST", "/save", reqBody)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), "user_id", 1.0)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()
	testHandler.saveURL(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected: %v, got: %v", http.StatusBadRequest, recorder.Code)
	}
}

func TestSaveURL_InvalidTime(t *testing.T){
	ms := &mocks.MockStorage{
		SetPairFunc: func(key, value string, userID int, expiresAt *time.Time) error {return service.ErrInvalidTime},
	}
	ls := &service.LinkService{
		Storage: ms,
		BaseURL: "http://localhost:8080",
	}
	testHandler := Handler{LinkService: ls}

	reqBody := strings.NewReader(`{"url": "https://google.com", "alias": "test", "duration":"1h"}`)
	req := httptest.NewRequest("POST", "/save", reqBody)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), "user_id", 1.0)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()
	testHandler.saveURL(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected: %v, got: %v", http.StatusBadRequest, recorder.Code)
	}
}

func TestSaveURL_Conflict(t *testing.T){
	ms := &mocks.MockStorage{
		SetPairFunc: func(key, value string, userID int, expiresAt *time.Time) error {return service.ErrAliasAlreadyExists},
	}
	ls := &service.LinkService{
		Storage: ms,
		BaseURL: "http://localhost:8080",
	}
	testHandler := Handler{LinkService: ls}

	reqBody := strings.NewReader(`{"url": "https://google.com", "alias": "test", "duration":"1h"}`)
	req := httptest.NewRequest("POST", "/save", reqBody)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), "user_id", 1.0)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()
	testHandler.saveURL(recorder, req)

	if recorder.Code != http.StatusConflict {
		t.Errorf("Expected: %v, got: %v", http.StatusConflict, recorder.Code)
	}
}

func TestSaveURL_InternalError(t *testing.T) {
	ms := &mocks.MockStorage{
		SetPairFunc: func(key, value string, userID int, expiresAt *time.Time) error {return errors.New("database server error")},
	}
	ls := &service.LinkService{
		Storage: ms,
		BaseURL: "http://localhost:8080",
	}
	testHandler := Handler{LinkService: ls}

	reqBody := strings.NewReader(`{"url": "https://google.com", "alias": "test", "duration":"1h"}`)
	req := httptest.NewRequest("POST", "/save", reqBody)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), "user_id", 1.0)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()
	testHandler.saveURL(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected: %v, got: %v", http.StatusInternalServerError, recorder.Code)
	}
}

func TestRedirectURL_Success(t *testing.T) {
	futureTime := time.Now().Add(time.Hour * 999)
	ms := &mocks.MockStorage{
		GetPairFunc: func(path string) (string, *time.Time, error) {return "https://google.com", &futureTime, nil},
	}
	ls := &service.LinkService{
		Storage: ms,
		BaseURL: "http://localhost:8080",
	}
	testHandler := Handler{LinkService: ls}

	req := httptest.NewRequest("GET", "/testAlias", nil)

	recorder := httptest.NewRecorder()
	testHandler.redirectURL(recorder, req)

	if recorder.Code != http.StatusFound {
		t.Errorf("Expected: %d, got: %d", http.StatusFound, recorder.Code)
	}

	location := recorder.Header().Get("Location")
	if location != "https://google.com" {
		t.Errorf("Expected: %s, got: %s", "https://google.com", location)
	}
}
func TestRedirectURL_InternalError(t *testing.T) {
	ms := &mocks.MockStorage{
		GetPairFunc: func(path string) (string, *time.Time, error) {return "", nil, errors.New("database server error")},
	}
	ls := &service.LinkService{
		Storage: ms,
		BaseURL: "http://localhost:8080",
	}
	testHandler := Handler{LinkService: ls}

	req := httptest.NewRequest("GET", "/testAlias", nil)

	recorder := httptest.NewRecorder()
	testHandler.redirectURL(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected: %d, got: %d", http.StatusInternalServerError, recorder.Code)
	}
}

func TestRedirectURL_NotFound(t *testing.T) {
	ms := &mocks.MockStorage{
		GetPairFunc: func(path string) (string, *time.Time, error) {return "", nil, store.ErrNotFound},
	}
	ls := &service.LinkService{
		Storage: ms,
		BaseURL: "http://localhost:8080",
	}
	testHandler := Handler{LinkService: ls}

	req := httptest.NewRequest("GET", "/testAlias", nil)

	recorder := httptest.NewRecorder()
	testHandler.redirectURL(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected: %d, got: %d", http.StatusNotFound, recorder.Code)
	}
}

func TestRedirectURL_TimeExpired(t *testing.T) {
	pastTime := time.Now().Add(time.Hour * -999)
	ms := &mocks.MockStorage{
		GetPairFunc: func(path string) (string, *time.Time, error) {return "https://google.com", &pastTime, nil},
	}
	ls := &service.LinkService{
		Storage: ms,
		BaseURL: "http://localhost:8080",
	}
	testHandler := Handler{LinkService: ls}

	req := httptest.NewRequest("GET", "/testAlias", nil)

	recorder := httptest.NewRecorder()
	testHandler.redirectURL(recorder, req)

	if recorder.Code != http.StatusGone {
		t.Errorf("Expected: %d, got: %d", http.StatusGone, recorder.Code)
	}

}

func TestCreateUser_Success(t *testing.T) {
	ms := &mocks.MockStorage{
		CreateUserFunc: func(username, password string) error {return nil},
	}
	as := &service.AuthService{
		Storage: ms,
	}
	testHandler := Handler{AuthService: as}

	reqBody := strings.NewReader(`{"username": "test", "password": "testPassword"}`)
	req := httptest.NewRequest("POST", "/register", reqBody)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	testHandler.createUser(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected: %d, got: %d", http.StatusCreated, recorder.Code)
	}
}

func TestCreateUser_InternalError(t *testing.T) {
	ms := &mocks.MockStorage{
		CreateUserFunc: func(username, password string) error {return errors.New("database server error")},
	}
	as := &service.AuthService{
		Storage: ms,
	}
	testHandler := Handler{AuthService: as}
	reqBody := strings.NewReader(`{"username": "test", "password": "testPassword"}`)
	req := httptest.NewRequest("POST", "/register", reqBody)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	testHandler.createUser(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected: %d, got: %d", http.StatusInternalServerError, recorder.Code)
	}
}

func TestCreateUser_Conflict(t *testing.T) {
	ms := &mocks.MockStorage{
		CreateUserFunc: func(username, password string) error {return service.ErrUserAlreadyExists},
	}
	as := &service.AuthService{
		Storage: ms,
	}
	testHandler := Handler{AuthService: as}
	reqBody := strings.NewReader(`{"username": "test", "password": "testPassword"}`)
	req := httptest.NewRequest("POST", "/register", reqBody)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	testHandler.createUser(recorder, req)

	if recorder.Code != http.StatusConflict {
		t.Errorf("Expected: %d, got: %d", http.StatusConflict, recorder.Code)
	}
}

func TestCreateUser_BadRequest(t *testing.T) {
	ms := &mocks.MockStorage{
		CreateUserFunc: func(username, password string) error {return nil},
	}
	as := &service.AuthService{
		Storage: ms,
	}
	testHandler := Handler{AuthService: as}
	reqBody := strings.NewReader(`{"username": "test", "password": ""}`)
	req := httptest.NewRequest("POST", "/register", reqBody)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	testHandler.createUser(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Expected: %d, got: %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	hashedPasswd, _ := bcrypt.GenerateFromPassword([]byte("testPassword"), bcrypt.DefaultCost)
	ms := &mocks.MockStorage{
		GetUserFunc: func(username string) (int, string, error) {return 0, string(hashedPasswd), nil},
	}
	as := &service.AuthService{
		Storage: ms,
	}
	testHandler := Handler{AuthService: as}
	reqBody := strings.NewReader(`{"username": "testUser", "password": "testPassword"}`)
	req := httptest.NewRequest("POST", "/login", reqBody)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	testHandler.login(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected: %d, got: %d", http.StatusOK, recorder.Code)
	}
}

func TestLogin_InternalError(t *testing.T) {
	ms := &mocks.MockStorage{
		GetUserFunc: func(username string) (int, string, error) {return 0, "", errors.New("database server error")},
	}
	as := &service.AuthService{
		Storage: ms,
	}
	testHandler := Handler{AuthService: as}
	reqBody := strings.NewReader(`{"username": "testUser", "password": "testPassword"}`)
	req := httptest.NewRequest("POST", "/login", reqBody)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	testHandler.login(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected: %d, got: %d", http.StatusInternalServerError, recorder.Code)
	}
}

func TestLogin_Unauthorized(t *testing.T) {
	ms := &mocks.MockStorage{
		GetUserFunc: func(username string) (int, string, error) {return 0, "", store.ErrNotFound},
	}
	as := &service.AuthService{
		Storage: ms,
	}
	testHandler := Handler{AuthService: as}
	reqBody := strings.NewReader(`{"username": "testUser", "password": "testPassword"}`)
	req := httptest.NewRequest("POST", "/login", reqBody)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	testHandler.login(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Expected: %d, got: %d", http.StatusUnauthorized, recorder.Code)
	}
}
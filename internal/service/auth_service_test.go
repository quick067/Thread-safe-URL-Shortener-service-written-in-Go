package service

import (
	"errors"
	"fmt"
	"testing"

	"URL-shortener/internal/store"
	"URL-shortener/internal/mocks"

	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser_Success(t *testing.T) {
	ms := &mocks.MockStorage{
		CreateUserFunc: func(username, password string) error {
			if password != "secretTestPassword" {
				fmt.Println("password hashed")
			}
			return nil
		},
	}
	as := &AuthService{
		Storage: ms,
	}

	err := as.CreateUser("testUser", "secretTestPassword")

	if err != nil {
		t.Errorf("Expected: no errors, got: %v", err)
	}
}

func TestCreateUser_InvalidPasswd(t *testing.T) {
	ms := &mocks.MockStorage{
		CreateUserFunc: func(username, password string) error {return nil},
	}
	as := &AuthService{
		Storage: ms,
	}

	err := as.CreateUser("testUser", "")

	if err == nil {
		t.Errorf("Expected: %v, got: %v", ErrInvalidPassword, err)
	}
	if !errors.Is(err, ErrInvalidPassword){
		t.Errorf("Expected: %v, got: %v", ErrInvalidPassword, err)
	}
}

func TestCreateUser_Conflict(t *testing.T) {
	ms := &mocks.MockStorage{
		CreateUserFunc: func(username, password string) error {return store.ErrConflict},
	}
	as := &AuthService{
		Storage: ms,
	}

	err := as.CreateUser("testUser", "secterTestPassword")

	if err == nil {
		t.Errorf("Expected: %v, got: %v", ErrUserAlreadyExists, err)
	}
	if !errors.Is(err, ErrUserAlreadyExists){
		t.Errorf("Expected: %v, got: %v", ErrUserAlreadyExists, err)
	}
}

func TestLogin_Success(t *testing.T) {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("testPassword"), bcrypt.DefaultCost)
	ms := &mocks.MockStorage{
		GetUserFunc: func(username string) (int, string, error) {return 1, string(hashedPass), nil},
	}
	as := &AuthService{
		Storage: ms,
	}

	token, err := as.LoginUser("testUser", "testPassword")

	if err != nil {
		t.Errorf("Expected: no errors, got: %v", err)
	}

	if token == "" {
		t.Errorf("Expected: non empty string, got: %s", token)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	ms := &mocks.MockStorage{
		GetUserFunc: func(username string) (int, string, error) {return 0, "", store.ErrNotFound},
	}
	as := &AuthService{
		Storage: ms,
	}
	_, err := as.LoginUser("testUser", "testPassword")

	if err == nil {
		t.Errorf("Expected: %v, got: %v", ErrInvalidCredentials, err)
	}
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Expected: %v, got: %v", ErrInvalidCredentials, err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("testPassword"), bcrypt.DefaultCost)
	ms := &mocks.MockStorage{
		GetUserFunc: func(username string) (int, string, error) {return 1, string(hashedPass), nil},
	}
	as := &AuthService{
		Storage: ms,
	}
	_, err := as.LoginUser("testUser", "secretPassword")

	if err == nil {
		t.Errorf("Expected: %v, got: %v", ErrInvalidCredentials, err)
	}

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Expected: %v, got: %v", ErrInvalidCredentials, err)
	}
}
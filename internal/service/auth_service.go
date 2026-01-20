package service

import (
	"errors"
	"fmt"
	"time"

	"URL-shortener/internal/store"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	Storage store.Storage
	JWTSecret string
}

func (as *AuthService) CreateUser(username, password string) error {
	if password == "" || len(password) < 8 {
		return fmt.Errorf("%w", ErrInvalidPassword)
	}

	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing passwd")
	}

	if err := as.Storage.CreateUser(username, string(hashedPasswd)); err != nil {
		if errors.Is(err, store.ErrConflict){
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (as *AuthService) LoginUser(username, password string) (string, error){
	id, hashedPasswd, err := as.Storage.GetUser(username)
	if err != nil {
		if errors.Is(err, store.ErrNotFound){
			return "", fmt.Errorf("%w, %s", ErrInvalidCredentials, store.ErrNotFound)
		}
		return "", err
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(password)); err != nil{
		return "", fmt.Errorf("%w", ErrInvalidCredentials)
	}

	claims := jwt.MapClaims{"user_id": id, "exp": time.Now().Add(time.Hour * 2).Unix()}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := newToken.SignedString([]byte(as.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("Error generating string token")
	}
	return tokenString, nil
}
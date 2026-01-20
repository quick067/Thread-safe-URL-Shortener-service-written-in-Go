package store

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("data conflict")
)


type Storage interface {
	GetPair(path string) (string, *time.Time, error)
	SetPair(key, value string, userID int, expiresAt *time.Time) error
	CreateUser(username, password string) error
	GetUser(username string) (int, string, error)
	Close() error
}


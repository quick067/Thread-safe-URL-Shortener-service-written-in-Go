package service

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"URL-shortener/internal/store"
)

var (
	ErrInvalidURL         = errors.New("invalid URL format")
	ErrInvalidTime        = errors.New("invalid duration")
	ErrAliasAlreadyExists = errors.New("alias already exists")
	ErrLinkNotFound       = errors.New("link not found")
	ErrLinkExpired        = errors.New("link time expired")
)

type LinkService struct {
	Storage store.Storage
	BaseURL string
}

func (ls *LinkService) CreateLink(originalURL, alias, duration string, userID int) (string, error) {
	if err := validateURL(originalURL); err != nil {
		return "", err
	}

	expiresAt, err := calcExpireTime(duration)
	if err != nil {
		return "", err
	}
	var URLKey string
	if alias != "" {
		URLKey = alias
		err := ls.saveCustomKey(URLKey, originalURL, userID, expiresAt)
		if err != nil {
			return "", err
		}
	} else {
		URLKey, err = ls.saveRandomKey(originalURL, userID, expiresAt)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s/%s", ls.BaseURL, URLKey), nil
}

func validateURL(URL string) error {
	parsedURL, err := url.ParseRequestURI(URL)
	if err != nil || URL == ""{
		return fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("%w: inadmissible protocol", ErrInvalidURL)
	}
	return nil
}

func (ls *LinkService) saveCustomKey(URLKey, originalURL string, userID int, expiresAt *time.Time) error {
	err := ls.Storage.SetPair(URLKey, originalURL, userID, expiresAt)
	if errors.Is(err, store.ErrConflict) {
		return fmt.Errorf("%w, %s", err, ErrAliasAlreadyExists)
	}
	return err
}

func (ls *LinkService) saveRandomKey(originalURL string, userID int, expiresAt *time.Time) (string, error) {
	for i := 0; i < 10; i++ {
		key := keyGenerator()

		err := ls.Storage.SetPair(key, originalURL, userID, expiresAt)
		if err == nil {
			return key, nil
		}

		if errors.Is(err, store.ErrConflict) {
			continue
		}

		return "", err
	}
	return "", fmt.Errorf("Error generating unique keys")
}

func calcExpireTime(duration string) (*time.Time, error) {
	var expiresAt *time.Time
	if duration != "" {
		parsedTime, err := time.ParseDuration(duration)
		if err != nil {
			return nil, ErrInvalidTime
		}
		expiration := time.Now().UTC().Add(parsedTime)
		expiresAt = &expiration
	}
	return expiresAt, nil
}

func (ls *LinkService) GetURL(alias string) (string, error) {
	if alias == "" {
		return "", errors.New("alias cannot be empty")
	}
	originalURL, expiresAt, err := ls.Storage.GetPair(alias)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return "", fmt.Errorf("%w, %w", err, ErrLinkNotFound)
		}
		return "", err
	}

	if expiresAt != nil && time.Now().UTC().After(*expiresAt) {
		return "", fmt.Errorf("%w", ErrLinkExpired)
	}
	return originalURL, nil
}

package store

import (
	"encoding/json"
	"os"
	"sync"
)

type Store struct {
	urlStore map[string]string
	mutex    sync.RWMutex
	filename string
}

func New(filename string) *Store{
	return &Store{
		urlStore: make(map[string]string),
		filename: filename,
	}
}

func (s *Store) SaveToFile() error {
	data, err := json.MarshalIndent(s.urlStore, "", " ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(s.filename, data, 0644); err != nil {
		return err
	}
	return nil
}

func (s *Store) LoadFromFile() error {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, &s.urlStore); err != nil {
		return err
	}
	return nil
}

func (s *Store) SetPair(key, value string) error {
	s.mutex.Lock()
	s.urlStore[key] = value
	err := s.SaveToFile()
	s.mutex.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetPair(path string) (string, bool) {
	s.mutex.RLock()
	value, ok := s.urlStore[path]
	s.mutex.RUnlock()
	return value, ok
}

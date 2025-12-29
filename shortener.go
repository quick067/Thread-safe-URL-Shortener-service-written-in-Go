package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type URLShortener struct {
	urlStore map[string]string
	mutex    sync.RWMutex
}

const fileName = "storage.json"

func (URLs *URLShortener) saveToFile() error {
	data, err := json.MarshalIndent(URLs.urlStore, "", " ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(fileName, data, 0644); err != nil {
		return err
	}
	return nil
}

func (URLs *URLShortener) loadFromFile() error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err = json.Unmarshal(data, &URLs.urlStore); err != nil {
		fmt.Println("Cannot parse JSON file")
		return err
	}
	return nil
}

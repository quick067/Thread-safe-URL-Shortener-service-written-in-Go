package mocks

import(
	"time"
)

type MockStorage struct {
	GetPairFunc func(path string) (string, *time.Time, error)
	SetPairFunc func(key, value string, userID int, expiresAt *time.Time) error
	CreateUserFunc func(username, password string) error
	GetUserFunc func(username string) (int, string, error)
	CloseFunc func() error
}

func (ms *MockStorage) GetPair(path string) (string, *time.Time, error) {
	if ms.GetPairFunc != nil {
		return ms.GetPairFunc(path)
	}
	return "", nil, nil 
}


func (ms *MockStorage) SetPair(key, value string, userID int, expiresAt *time.Time) error {
	if ms.SetPairFunc != nil {
		return ms.SetPairFunc(key, value, userID, expiresAt)
	}
	return nil 
}

func (ms *MockStorage) CreateUser(username, password string) error {
	if ms.CreateUserFunc != nil {
		return ms.CreateUserFunc(username, password)
	}
	return nil
}

func (ms *MockStorage) GetUser(username string) (int, string, error) {
	if ms.GetUserFunc != nil {
		return ms.GetUserFunc(username)
	}
	return 0, "", nil 
}

func (ms *MockStorage) Close() error {
	if ms.CloseFunc != nil {
		return ms.CloseFunc()
	}
	return nil
}
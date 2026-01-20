package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"


	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func InitDB(DSN string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return &PostgresStore{}, fmt.Errorf("error opening DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return &PostgresStore{}, fmt.Errorf("error pinging DB: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (ps *PostgresStore) Close() error {
	if err := ps.db.Close(); err != nil {
		return fmt.Errorf("error closing db: %w", err)
	}
	return nil
}

func (ps *PostgresStore) CreateKeysTable() error {
	sqlQuery := `CREATE TABLE IF NOT EXISTS links(
					short_key TEXT PRIMARY KEY,
					original_key TEXT,
					user_id INTEGER REFERENCES users(id),
					expires_at TIMESTAMP
	);`

	if _, err := ps.db.Exec(sqlQuery); err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	return nil
}

func (ps *PostgresStore) CreateUsersTable() error {
	sqlQuery := `CREATE TABLE IF NOT EXISTS users(
					id SERIAL PRIMARY KEY,
					username TEXT NOT NULL UNIQUE,
					hashed_passwd TEXT NOT NULL				
	)`
	if _, err := ps.db.Exec(sqlQuery); err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	return nil
}

func (ps *PostgresStore) CreateTables() error {
	if err := ps.CreateKeysTable(); err != nil {
		return err
	}

	if err := ps.CreateUsersTable(); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) SetPair(key, value string, userID int, expiresAt *time.Time) error {
	sqlQuery := `INSERT INTO links (short_key, original_key, user_id, expires_at)
					VALUES ($1, $2, $3, $4);`

	_, err := ps.db.Exec(sqlQuery, key, value, userID, expiresAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr){
			if pqErr.Code == "23505"{
				return ErrConflict
			}
		}
		return err
	}
	return nil
}

func (ps *PostgresStore) GetPair(path string) (string, *time.Time, error) {
	var result string
	var expiresAt *time.Time
	sqlQuery := `SELECT original_key, expires_at FROM links WHERE short_key = $1`
	err := ps.db.QueryRow(sqlQuery, path).Scan(&result, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, ErrNotFound
		}
		return "", nil, err
	}

	return result, expiresAt, nil
}

func (ps *PostgresStore) CreateUser(username, hashedPass string) error {
	sqlQuery := `INSERT INTO users (username, hashed_passwd)
					VALUES ($1, $2);`
	if _, err := ps.db.Exec(sqlQuery, username, string(hashedPass)); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr){
			if pqErr.Code == "23505"{
				return ErrConflict
			}
		}
		return fmt.Errorf("error inserting user info: %w", err)
	}
	return nil
}

func (ps *PostgresStore) GetUser(username string) (int, string, error) {
	sqlQuery := `SELECT id, hashed_passwd FROM users WHERE username = $1`
	var id int
	var hashed_passwd string

	if err := ps.db.QueryRow(sqlQuery, username).Scan(&id, &hashed_passwd); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", ErrNotFound
		}
		return 0, "", err
	}
	return id, hashed_passwd, nil
}

package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"myURLShortener/internal/http-server/handlers/url/delete"
	"myURLShortener/internal/http-server/handlers/url/showall"
	"myURLShortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url,alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert ID: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) GetAliasAndURL(alias, url string) (string, string, error) {
	const op = "storage.sqlite.GetAliasAndURL"

	stmt, err := s.db.Prepare("SELECT alias, url FROM url WHERE alias = ? AND url = ?")
	if err != nil {
		return "", "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resAlias, resURL string

	err = stmt.QueryRow(alias, url).Scan(&resAlias, &resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", storage.ErrURLNotFound
		}

		return "", "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resAlias, resURL, nil
}

func (s *Storage) DeleteByAliasAndURL(alias, url string) error {
	const op = "storage.sqlite.DeleteByAliasAndURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ? AND url = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(alias, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrURLNotFound
		}

		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteURLByAlias(alias string) error {
	const op = "storage.sqlite.DeleteURLByAlias"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrURLNotFound
		}

		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteAliasByURL(url string) ([]*delete.AliasData, error) {
	const op = "storage.sqlite.DeleteAliasByURL"

	stmtOne, err := s.db.Prepare("SELECT alias FROM url WHERE url = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	rows, err := stmtOne.Query(url)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	aliasData := make([]*delete.AliasData, 0)

	for rows.Next() {
		item := new(delete.AliasData)
		err = rows.Scan(&item.Alias)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		aliasData = append(aliasData, item)
	}

	stmtTwo, err := s.db.Prepare("DELETE FROM url WHERE url = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmtTwo.Exec(url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrURLNotFound
		}

		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return aliasData, nil
}

func (s *Storage) ShowURL() ([]*showall.AliasUrl, error) {
	const op = "storage.sqlite.ShowURL"

	stmt, err := s.db.Prepare("SELECT alias, url FROM url")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	aliasUrl := make([]*showall.AliasUrl, 0)

	for rows.Next() {
		pair := new(showall.AliasUrl)
		err = rows.Scan(&pair.Alias, &pair.URL)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		aliasUrl = append(aliasUrl, pair)
	}

	return aliasUrl, nil
}

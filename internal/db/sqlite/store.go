package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stageddat/shelter-node/internal/entry"
	"github.com/stageddat/shelter-node/internal/user"
)

type SQLiteStore struct {
	db *sql.DB
}

func New(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("sqlite ping: %w", err)
	}
	s := &SQLiteStore{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("sqlite migrate: %w", err)
	}
	return s, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			display_name TEXT NOT NULL DEFAULT '',
			auth_key_hash TEXT NOT NULL,
			encrypted_master_key BLOB,
			salt BLOB,
			iv BLOB,
			recovery_encrypted_master_key BLOB,
			recovery_salt BLOB,
			recovery_iv BLOB,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS entries (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			encrypted_title BLOB,
			title_iv BLOB,
			encrypted_content BLOB,
			content_iv BLOB,
			date TEXT,
			time TEXT,
			word_count INTEGER DEFAULT 0,
			char_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

// --- user ---

func (s *SQLiteStore) CreateUser(ctx context.Context, u *user.User) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users 
			(id, username, display_name, auth_key_hash, encrypted_master_key, salt, iv,
			recovery_encrypted_master_key, recovery_salt, recovery_iv)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.DisplayName, u.AuthKeyHash,
		u.EncryptedMasterKey, u.Salt, u.IV,
		u.RecoveryEncryptedMasterKey, u.RecoverySalt, u.RecoveryIV,
	)
	return err
}

func (s *SQLiteStore) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, username, display_name, auth_key_hash, encrypted_master_key, salt, iv,
			recovery_encrypted_master_key, recovery_salt, recovery_iv,
			created_at, updated_at
		FROM users WHERE username = ?`,
		username,
	)
	u := &user.User{}
	err := row.Scan(
		&u.ID, &u.Username, &u.DisplayName, &u.AuthKeyHash,
		&u.EncryptedMasterKey, &u.Salt, &u.IV,
		&u.RecoveryEncryptedMasterKey, &u.RecoverySalt, &u.RecoveryIV,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return u, err
}

// --- entry ---

func (s *SQLiteStore) CreateEntry(ctx context.Context, e *entry.Entry) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO entries
			(id, user_id, encrypted_title, title_iv, encrypted_content, content_iv,
			date, time, word_count, char_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.UserID, e.EncryptedTitle, e.TitleIV,
		e.EncryptedContent, e.ContentIV,
		e.Date, e.Time, e.WordCount, e.CharCount,
	)
	return err
}

func (s *SQLiteStore) GetEntry(ctx context.Context, id string) (*entry.Entry, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, user_id, encrypted_title, title_iv, encrypted_content, content_iv,
			date, time, word_count, char_count, created_at, updated_at
		FROM entries WHERE id = ?`,
		id,
	)
	e := &entry.Entry{}
	err := row.Scan(
		&e.ID, &e.UserID, &e.EncryptedTitle, &e.TitleIV,
		&e.EncryptedContent, &e.ContentIV,
		&e.Date, &e.Time, &e.WordCount, &e.CharCount,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("entry not found")
	}
	return e, err
}

func (s *SQLiteStore) GetEntriesByUser(ctx context.Context, userID string) ([]*entry.Entry, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, user_id, encrypted_title, title_iv, encrypted_content, content_iv,
			date, time, word_count, char_count, created_at, updated_at
		FROM entries WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*entry.Entry
	for rows.Next() {
		e := &entry.Entry{}
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.EncryptedTitle, &e.TitleIV,
			&e.EncryptedContent, &e.ContentIV,
			&e.Date, &e.Time, &e.WordCount, &e.CharCount,
			&e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (s *SQLiteStore) UpdateEntry(ctx context.Context, e *entry.Entry) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE entries SET
			encrypted_title = ?, title_iv = ?,
			encrypted_content = ?, content_iv = ?,
			word_count = ?, char_count = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		e.EncryptedTitle, e.TitleIV,
		e.EncryptedContent, e.ContentIV,
		e.WordCount, e.CharCount,
		e.ID,
	)
	return err
}

func (s *SQLiteStore) DeleteEntry(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM entries WHERE id = ?`, id)
	return err
}

package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/stageddat/shelter-node/internal/entry"
	"github.com/stageddat/shelter-node/internal/user"
)

type PostgresStore struct {
	db *pgx.Conn
}

func New(connString string) (*PostgresStore, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}
	s := &PostgresStore{db: conn}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("postgres migrate: %w", err)
	}
	return s, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close(context.Background())
}

func (s *PostgresStore) migrate() error {
	_, err := s.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			display_name TEXT NOT NULL DEFAULT '',
			auth_key_hash TEXT NOT NULL,
			encrypted_master_key BYTEA,
			salt BYTEA,
			iv BYTEA,
			recovery_encrypted_master_key BYTEA,
			recovery_salt BYTEA,
			recovery_iv BYTEA,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS entries (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			encrypted_title BYTEA,
			title_iv BYTEA,
			encrypted_content BYTEA,
			content_iv BYTEA,
			date TEXT,
			time TEXT,
			word_count INTEGER DEFAULT 0,
			char_count INTEGER DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
	`)
	return err
}

// --- user ---

func (s *PostgresStore) CreateUser(ctx context.Context, u *user.User) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO users
			(id, username, display_name, auth_key_hash, encrypted_master_key, salt, iv,
			recovery_encrypted_master_key, recovery_salt, recovery_iv)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		u.ID, u.Username, u.DisplayName, u.AuthKeyHash,
		u.EncryptedMasterKey, u.Salt, u.IV,
		u.RecoveryEncryptedMasterKey, u.RecoverySalt, u.RecoveryIV,
	)
	return err
}

func (s *PostgresStore) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	row := s.db.QueryRow(ctx,
		`SELECT id, username, display_name, auth_key_hash, encrypted_master_key, salt, iv,
			recovery_encrypted_master_key, recovery_salt, recovery_iv,
			created_at, updated_at
		FROM users WHERE username = $1`,
		username,
	)
	u := &user.User{}
	err := row.Scan(
		&u.ID, &u.Username, &u.DisplayName, &u.AuthKeyHash,
		&u.EncryptedMasterKey, &u.Salt, &u.IV,
		&u.RecoveryEncryptedMasterKey, &u.RecoverySalt, &u.RecoveryIV,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return u, err
}

// --- entry ---

func (s *PostgresStore) CreateEntry(ctx context.Context, e *entry.Entry) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO entries
			(id, user_id, encrypted_title, title_iv, encrypted_content, content_iv,
			date, time, word_count, char_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		e.ID, e.UserID, e.EncryptedTitle, e.TitleIV,
		e.EncryptedContent, e.ContentIV,
		e.Date, e.Time, e.WordCount, e.CharCount,
	)
	return err
}

func (s *PostgresStore) GetEntry(ctx context.Context, id string) (*entry.Entry, error) {
	row := s.db.QueryRow(ctx,
		`SELECT id, user_id, encrypted_title, title_iv, encrypted_content, content_iv,
			date, time, word_count, char_count, created_at, updated_at
		FROM entries WHERE id = $1`,
		id,
	)
	e := &entry.Entry{}
	err := row.Scan(
		&e.ID, &e.UserID, &e.EncryptedTitle, &e.TitleIV,
		&e.EncryptedContent, &e.ContentIV,
		&e.Date, &e.Time, &e.WordCount, &e.CharCount,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("entry not found")
	}
	return e, err
}

func (s *PostgresStore) GetEntriesByUser(ctx context.Context, userID string) ([]*entry.Entry, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, encrypted_title, title_iv, encrypted_content, content_iv,
			date, time, word_count, char_count, created_at, updated_at
		FROM entries WHERE user_id = $1 ORDER BY created_at DESC`,
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

func (s *PostgresStore) UpdateEntry(ctx context.Context, e *entry.Entry) error {
	_, err := s.db.Exec(ctx,
		`UPDATE entries SET
			encrypted_title = $1, title_iv = $2,
			encrypted_content = $3, content_iv = $4,
			word_count = $5, char_count = $6,
			updated_at = NOW()
		WHERE id = $7`,
		e.EncryptedTitle, e.TitleIV,
		e.EncryptedContent, e.ContentIV,
		e.WordCount, e.CharCount,
		e.ID,
	)
	return err
}

func (s *PostgresStore) DeleteEntry(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM entries WHERE id = $1`, id)
	return err
}

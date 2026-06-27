package db

import (
	"context"

	"github.com/stageddat/shelter-node/internal/entry"
	"github.com/stageddat/shelter-node/internal/user"
)

type Store interface {
	CreateUser(ctx context.Context, u *user.User) error
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)

	CreateEntry(ctx context.Context, e *entry.Entry) error
	GetEntry(ctx context.Context, id string) (*entry.Entry, error)
	GetEntriesByUser(ctx context.Context, userID string) ([]*entry.Entry, error)
	UpdateEntry(ctx context.Context, e *entry.Entry) error
	DeleteEntry(ctx context.Context, id string) error

	Close() error
}

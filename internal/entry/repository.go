package entry

import "context"

type Repository interface {
	CreateEntry(ctx context.Context, e *Entry) error
	GetEntry(ctx context.Context, id string) (*Entry, error)
	GetEntriesByUser(ctx context.Context, userID string) ([]*Entry, error)
	UpdateEntry(ctx context.Context, e *Entry) error
	DeleteEntry(ctx context.Context, id string) error
}

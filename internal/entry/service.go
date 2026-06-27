package entry

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateEntryRequest) (*Entry, error) {
	e := &Entry{
		ID:               uuid.NewString(),
		UserID:           req.UserID,
		EncryptedTitle:   req.EncryptedTitle,
		TitleIV:          req.TitleIV,
		EncryptedContent: req.EncryptedContent,
		ContentIV:        req.ContentIV,
		Date:             req.Date,
		Time:             req.Time,
		WordCount:        req.WordCount,
		CharCount:        req.CharCount,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.CreateEntry(ctx, e); err != nil {
		return nil, fmt.Errorf("create entry: %w", err)
	}
	return e, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Entry, error) {
	e, err := s.repo.GetEntry(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get entry: %w", err)
	}
	return e, nil
}

func (s *Service) GetByUser(ctx context.Context, userID string) ([]*Entry, error) {
	entries, err := s.repo.GetEntriesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get entries: %w", err)
	}
	return entries, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateEntryRequest) (*Entry, error) {
	e, err := s.repo.GetEntry(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("entry not found: %w", err)
	}

	if req.EncryptedTitle != nil {
		e.EncryptedTitle = req.EncryptedTitle
		e.TitleIV = req.TitleIV
	}
	if req.EncryptedContent != nil {
		e.EncryptedContent = req.EncryptedContent
		e.ContentIV = req.ContentIV
	}
	e.UpdatedAt = time.Now()

	if err := s.repo.UpdateEntry(ctx, e); err != nil {
		return nil, fmt.Errorf("update entry: %w", err)
	}
	return e, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteEntry(ctx, id)
}

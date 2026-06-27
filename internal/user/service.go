package user

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	if req.Username == "" || req.AuthKeyHash == "" {
		return nil, fmt.Errorf("username and auth_key_hash are required")
	}

	u := &User{
		ID:          newID(),
		Username:    req.Username,
		DisplayName: req.Username,
		AuthKeyHash: req.AuthKeyHash,
	}

	if err := s.repo.CreateUser(ctx, u); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	u, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

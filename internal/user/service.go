package user

import (
	"context"
	"errors"
	"strings"
	"time"
)

var ErrInvalid = errors.New("некорректные данные пользователя")

type Service struct {
	repo Repository
	now  func() time.Time // подменяется в тестах
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// IsRegistered нужен для проверки в диалоге (типо вы уже зарегистрированны)
func (s *Service) IsRegistered(ctx context.Context, telegramID int64) (bool, error) {
	_, err := s.repo.Get(ctx, telegramID)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) Get(ctx context.Context, telegramID int64) (*User, error) {
	return s.repo.Get(ctx, telegramID)
}

// Register создаёт пользователя один раз
// Повторный вызов вернёт ErrAlreadyExists
// AI-Assisted поможет атомарность вместо проверки до
func (s *Service) Register(ctx context.Context, u User) (*User, error) {
	u.FullName = strings.TrimSpace(u.FullName)

	switch {
	case u.TelegramID == 0:
		return nil, ErrInvalid
	case u.FullName == "":
		return nil, ErrInvalid
	case u.ConsentAt.IsZero():
		return nil, ErrInvalid // без согласия на ПД регистрировать нельзя!!!
	}

	now := s.now()
	u.RegisteredAt = now
	u.UpdatedAt = now

	if err := s.repo.Create(ctx, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) List(ctx context.Context) ([]*User, error) {
	return s.repo.List(ctx)
}

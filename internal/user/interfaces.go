package user

import "context"

type Repository interface {
	// Get вернет копию пользователя или ErrNotFound
	Get(ctx context.Context, telegramID int64) (*User, error)

	// Create добавляет пользователя или возвращает ErrAlreadyExists
	// AI-Assisted - атомарность нормальная
	Create(ctx context.Context, u *User) error

	// List возвращает копии всех пользователей (для рассылок)
	List(ctx context.Context) ([]*User, error)
}

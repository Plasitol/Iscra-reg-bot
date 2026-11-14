package user

import (
	"errors"
	"time"
)

var (
	ErrAlreadyExists = errors.New("пользователь уже зарегистрирован")
	ErrNotFound      = errors.New("пользователь не найден")
)

type User struct {
	// требуемая информация
	TelegramID   int64  `json:"telegram_id"`
	FullName     string `json:"full_name"`
	Passport     int64  `json:"passport"`
	BirthDate    string `json:"birth_date"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Organization string `json:"organization"`

	// согласие на обработку ПД
	ConsentAt time.Time `json:"consent_at"`

	// служебное
	RegisteredAt time.Time `json:"registered_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

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
	// из Telegram
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username,omitempty"`

	// данные регистрации
	FullName     string `json:"full_name"`
	Passport     string `json:"passport"`
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

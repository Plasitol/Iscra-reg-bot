package bot

import (
	"github.com/go-telegram/bot/models"
)

// главное меню
func mainKeyboard() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "Регистрация"},
				{Text: "Помощь"},
			},
		},
		ResizeKeyboard: true,
		IsPersistent:   true,
	}
}

// согласие на обработку пд
func consentKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Согласен",
					CallbackData: "consent:accept",
				},
			},
			{
				{
					Text:         "Не согласен",
					CallbackData: "consent:decline",
				},
			},
		},
	}
}

// подтверждение реги
func confirmKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Подтвердить",
					CallbackData: "registration:confirm",
				},
			},
			{
				{
					Text:         "Изменить данные",
					CallbackData: "registration:edit",
				},
			},
		},
	}
}

// отмена реги
func cancelKeyboard() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "Отменить регистрацию"},
			},
		},
		ResizeKeyboard: true,
	}
}

// удаление клавиатуры
func removeKeyboard() *models.ReplyKeyboardRemove {
	return &models.ReplyKeyboardRemove{
		RemoveKeyboard: true,
	}
}

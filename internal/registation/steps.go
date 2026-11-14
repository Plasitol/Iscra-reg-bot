package registration

type Field struct {
	Key      string
	Prompt   string // ключ текста вопроса
	Error    string // ключ текста ошибки
	Validate func(string) (string, error)
}

var Fields = []Field{ // TODO допиши файлы
	{Key: "fio", Prompt: "texts/registration/questions/fio", Error: "texts/registration/errors/fio", Validate: validateFIO},
	{Key: "passport", Prompt: "texts/registration/questions/passport", Error: "texts/registration/errors/passport", Validate: validatePassport},
	{Key: "birth_date", Prompt: "texts/registration/questions/birth_date", Error: "texts/registration/errors/birth_date", Validate: validateBirthDate},
	{Key: "phone", Prompt: "texts/registration/questions/phone", Error: "texts/registration/errors/phone", Validate: validatePhone},
	{Key: "email", Prompt: "texts/registration/questions/email", Error: "texts/registration/errors/email", Validate: validateEmail},
	{Key: "organization", Prompt: "texts/registration/questions/organization", Error: "texts/registration/errors/organization", Validate: validateOrganization},
}

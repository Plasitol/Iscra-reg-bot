package registration

import (
	"errors"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// plasitol
// AI-Assisted

var ErrInvalidInput = errors.New("некорректный ввод")

// подменяется в тестах
var now = time.Now

const (
	minAge = 14
	maxAge = 100

	minNameWords = 2 // 3 если отчество обязтельно (?) TODO
	maxNameWords = 4
	maxNameLen   = 100

	maxEmailLen = 254
	maxOrgLen   = 150
)

// вводимые символы, проверка
var (
	fioRe = regexp.MustCompile(`^\p{L}[\p{L} '’-]*$`)

	passportRe = regexp.MustCompile(`^\d{10}$`)

	phoneAllowedRe = regexp.MustCompile(`^\+?[\d\s()\-]+$`)

	dateRe = regexp.MustCompile(`^(\d{1,2})[./-](\d{1,2})[./-](\d{4})$`)

	orgRe = regexp.MustCompile(`^[\p{L}\p{N}][\p{L}\p{N} .,'’"«»()№&/+\-]*$`)
)

// normalize обрезает края
func normalize(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// ФИО
func validateFIO(input string) (string, error) {
	s := normalize(input)

	if utf8.RuneCountInString(s) > maxNameLen || !fioRe.MatchString(s) {
		return "", ErrInvalidInput
	}
	if n := len(strings.Fields(s)); n < minNameWords || n > maxNameWords {
		return "", ErrInvalidInput
	}
	return s, nil
}

// Паспорт РФ (хранится нормально)
func validatePassport(input string) (string, error) {
	digits := strings.Join(strings.Fields(input), "")
	if !passportRe.MatchString(digits) {
		return "", ErrInvalidInput
	}

	series, number := digits[:4], digits[4:]
	if series == "0000" || number == "000000" {
		return "", ErrInvalidInput
	}
	return series + " " + number, nil
}

// Дата рождения
func validateBirthDate(input string) (string, error) {
	m := dateRe.FindStringSubmatch(strings.TrimSpace(input))
	if m == nil {
		return "", ErrInvalidInput
	}

	day, _ := strconv.Atoi(m[1])
	month, _ := strconv.Atoi(m[2])
	year, _ := strconv.Atoi(m[3])

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	// time.Date нормализует 31.02 в 03.03
	if t.Day() != day || int(t.Month()) != month {
		return "", ErrInvalidInput
	}

	today := now()
	age := today.Year() - year
	if today.Month() < t.Month() || (today.Month() == t.Month() && today.Day() < day) {
		age--
	}
	// дата в будущем даёт отрицательный возраст и тоже отсекается
	if age < minAge || age > maxAge {
		return "", ErrInvalidInput
	}
	return t.Format("02.01.2006"), nil
}

// Телефон РФ
func validatePhone(input string) (string, error) {
	s := strings.TrimSpace(input)
	if !phoneAllowedRe.MatchString(s) {
		return "", ErrInvalidInput
	}

	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	d := b.String()
	plus := strings.HasPrefix(s, "+")

	switch {
	case len(d) == 11 && d[0] == '7':
		d = d[1:]
	case len(d) == 11 && d[0] == '8' && !plus:
		d = d[1:]
	case len(d) == 10:
		// код страны опущен
	default:
		return "", ErrInvalidInput
	}

	if d[0] != '9' {
		return "", ErrInvalidInput
	}
	return "+7" + d, nil
}

// Email
func validateEmail(input string) (string, error) {
	s := strings.TrimSpace(input)
	if s == "" || len(s) > maxEmailLen {
		return "", ErrInvalidInput
	}

	addr, err := mail.ParseAddress(s)
	if err != nil || addr.Address != s {
		return "", ErrInvalidInput
	}

	at := strings.LastIndex(s, "@")
	if !strings.Contains(s[at+1:], ".") { // user@localhost не годится!!
		return "", ErrInvalidInput
	}
	return strings.ToLower(s), nil
}

// Вуз / Школа / Компания
func validateOrganization(input string) (string, error) {
	s := normalize(input)

	if n := utf8.RuneCountInString(s); n < 2 || n > maxOrgLen {
		return "", ErrInvalidInput
	}
	if !orgRe.MatchString(s) {
		return "", ErrInvalidInput
	}
	return s, nil
}

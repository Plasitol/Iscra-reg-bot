package messages

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"text/template"
)

// довольно опасно, не кладите в message env!!!

//go:embed all:*
var embedded embed.FS

// plasitol
// AI-Assisted, я довольно плох в работе с парсингом, однако нейронка помогла на ура

// Ключи как константы: опечатка станет ошибкой компиляции, а не рантайма.
const (
	// default-значения для вызова из кода
	Start   = "start"
	Help    = "default/help"
	Unknown = "default/unknown"

	// регистрация
	Registration = "registration/reg"
	Decline      = "registration/decline"
	Accept       = "registration/accept"
)

// Ключи, которые обязаны существовать (проверяются при старте)
var required = []string{Start, Help, Unknown, Accept, Decline, Registration}

type Store struct {
	tmpl *template.Template
}

// Load рекурсивно читает все *.md из fs.FS.
// Ключ = путь без расширения, например "default/help".
func Load(fsys fs.FS) (*Store, error) {
	t := template.New("")

	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".md") {
			return nil
		}

		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return fmt.Errorf("read %s: %w", p, err)
		}

		name := strings.TrimSuffix(p, ".md")
		if _, err := t.New(name).Parse(string(data)); err != nil {
			return fmt.Errorf("parse %s: %w", p, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Fail fast: проверяем при старте, а не когда юзер нажмёт /help.
	for _, key := range required {
		if t.Lookup(key) == nil {
			return nil, fmt.Errorf("отсутствует обязательный текст %q", key)
		}
	}

	return &Store{tmpl: t}, nil
}

// LoadEmbedded использует файлы, вшитые в бинарник.
func LoadEmbedded() (*Store, error) { return Load(embedded) }

// Get возвращает текст; data нужна для плейсхолдеров вида {{.Name}}.
func (s *Store) Get(key string, data any) (string, error) {
	t := s.tmpl.Lookup(key)
	if t == nil {
		return "", fmt.Errorf("текст %q не найден", key)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

package messages

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"text/template"
)

//go:embed *.md
var embedded embed.FS

// plasitol
// AI-Assisted, я довольно плох в работе с парсингом, однако нейронка помогла на ура

// Ключи как константы: опечатка станет ошибкой компиляции, а не рантайма.
const (
	// default-значения для вызова из кода
	Start   = "default/start"
	Help    = "default/help"
	Unknown = "default/unknown"
)

// Ключи, которые обязаны существовать.
var required = []string{Start, Help}

type Store struct {
	tmpl *template.Template
}

// Load читает все *.md из fs.FS. Ключ = имя файла без расширения.
func Load(fsys fs.FS) (*Store, error) {
	t := template.New("")

	files, err := fs.Glob(fsys, "*.md")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		data, err := fs.ReadFile(fsys, f)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", f, err)
		}
		name := strings.TrimSuffix(path.Base(f), ".md")
		if _, err := t.New(name).Parse(string(data)); err != nil {
			return nil, fmt.Errorf("parse %s: %w", f, err)
		}
	}

	s := &Store{tmpl: t}

	// Fail fast: проверяем при старте, а не когда юзер нажмёт /help.
	for _, key := range required {
		if t.Lookup(key) == nil {
			return nil, fmt.Errorf("отсутствует обязательный текст %q", key)
		}
	}
	return s, nil
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

package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"plasitol.dev/iscrav1/internal/user"
)

// plasitol
// AI-Assisted, я все еще плох в парсинге бля(

const fileVersion = 1

// формат файла - объект, а не голый массив (TODO миграция)
type fileFormat struct {
	Version int          `json:"version"`
	Users   []*user.User `json:"users"`
}

type JSONStore struct {
	mu    sync.RWMutex
	path  string
	users map[int64]*user.User
}

// проверка на этапе компиляции, что JSONStore реализует интерфейс
var _ user.Repository = (*JSONStore)(nil)

func NewJSON(path string) (*JSONStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("создание папки данных: %w", err)
	}

	s := &JSONStore{path: path, users: make(map[int64]*user.User)}

	data, err := os.ReadFile(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return s, nil // первый запуск, файла ещё нет
	case err != nil:
		return nil, fmt.Errorf("чтение %s: %w", path, err)
	case len(data) == 0:
		return s, nil
	}

	var f fileFormat
	if err := json.Unmarshal(data, &f); err != nil {
		// не затираем файл молча: пусть оператор разберётся
		return nil, fmt.Errorf("файл %s повреждён: %w", path, err)
	}
	for _, u := range f.Users {
		s.users[u.TelegramID] = u
	}
	return s, nil
}

func (s *JSONStore) Get(_ context.Context, id int64) (*user.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return nil, user.ErrNotFound
	}
	cp := *u // отдаём копию, чтобы снаружи не меняли данные в карте
	return &cp, nil
}

func (s *JSONStore) Create(_ context.Context, u *user.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[u.TelegramID]; ok {
		return user.ErrAlreadyExists
	}

	cp := *u
	s.users[u.TelegramID] = &cp

	if err := s.flushLocked(); err != nil {
		delete(s.users, u.TelegramID) // откат: в памяти не должно быть того, чего нет на диске
		return err
	}
	return nil
}

func (s *JSONStore) Update(_ context.Context, u *user.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	old, ok := s.users[u.TelegramID]
	if !ok {
		return user.ErrNotFound
	}

	cp := *u
	s.users[u.TelegramID] = &cp

	if err := s.flushLocked(); err != nil {
		s.users[u.TelegramID] = old // откат
		return err
	}
	return nil
}

func (s *JSONStore) List(_ context.Context) ([]*user.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*user.User, 0, len(s.users))
	for _, u := range s.users {
		cp := *u
		out = append(out, &cp)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].RegisteredAt.Before(out[j].RegisteredAt) })
	return out, nil
}

// flushLocked записывает всё на диск. Вызывать только под s.mu.Lock().
// Пишем во временный файл и переименовываем: на диске всегда либо старая, либо новая версия целиком.
func (s *JSONStore) flushLocked() error {
	list := make([]*user.User, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, u)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].TelegramID < list[j].TelegramID })

	data, err := json.MarshalIndent(fileFormat{Version: fileVersion, Users: list}, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(s.path), "users-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	cleanup := func() { _ = os.Remove(tmpName) }

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Sync(); err != nil { // данные реально на диске до rename
		tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return err
	}
	if err := os.Chmod(tmpName, 0o600); err != nil {
		cleanup()
		return err
	}
	if err := os.Rename(tmpName, s.path); err != nil {
		cleanup()
		return err
	}
	return nil
}

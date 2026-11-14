package registration

import (
	"maps"
	"sync"
	"time"
)

// TODO: очистика сессии, которая началсь раньше суток (зависшие сессии)

type Session struct {
	Step      int // индекс текущего вопроса
	Answers   map[string]string
	ConsentAt time.Time
	StartedAt time.Time
}

type Sessions struct {
	mu sync.Mutex
	m  map[int64]*Session
}

func NewSessions() *Sessions {
	return &Sessions{m: make(map[int64]*Session)}
}

// рестарт сессии при изменении данных
func (s *Sessions) Restart(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.m[id]; ok {
		sess.Step = 0
		sess.Answers = make(map[string]string)
	}
}

// Start новая сессия, но повторное согласие скинет её!!
func (s *Sessions) Start(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.m[id] = &Session{
		Answers:   make(map[string]string),
		ConsentAt: now,
		StartedAt: now,
	}
}

func (s *Sessions) Get(id int64) (Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.m[id]
	if !ok {
		return Session{}, false
	}
	cp := *sess
	cp.Answers = make(map[string]string, len(sess.Answers))
	maps.Copy(cp.Answers, sess.Answers)
	return cp, true
}

type Result int

const (
	Invalid  Result = iota // ответ не прошёл валидацию, остаёмся на шаге
	NextStep               // принято -> есть следующий вопрос
	Finished               // принят последний ответ -> пора показывать итог
	NoSession
)

// Submit валидирует ответ на текущем шаге и сдвигает сессию.
func (s *Sessions) Submit(id int64, input string) (Result, *Field) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.m[id]
	if !ok {
		return NoSession, nil
	}

	f := &Fields[sess.Step]
	value, err := f.Validate(input)
	if err != nil {
		return Invalid, f // удобно, подсказка ошибки
	}

	sess.Answers[f.Key] = value
	sess.Step++

	if sess.Step >= len(Fields) {
		return Finished, nil
	}
	return NextStep, &Fields[sess.Step]
}

func (s *Sessions) Delete(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, id)
}

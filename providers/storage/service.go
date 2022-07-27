package storage

import (
	"context"
	"github.com/cadyrov/occam/domain"
	"github.com/rs/zerolog"
	"sync"
)

type Service struct {
	log   *zerolog.Logger
	shift int

	mu      sync.RWMutex
	storage map[domain.Ticker]float64
}

func New(log *zerolog.Logger, shift int) *Service {
	return &Service{
		storage: make(map[domain.Ticker]float64),
		shift:   shift,
		log:     log,
	}
}

func (s *Service) Put(ctx context.Context, ti domain.Ticker, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil
	default:
		f := s.storage[ti]

		if f == 0 || s.shift < 2 {
			s.storage[ti] = value
		} else {
			f = f/float64(s.shift)*float64(s.shift-1) + value/float64(s.shift)
		}
	}

	return nil
}

func (s *Service) Get(ti domain.Ticker) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.storage[ti]
}

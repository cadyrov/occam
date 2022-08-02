package storage

import (
	"context"
	"github.com/cadyrov/occam/domain"
	"github.com/rs/zerolog"
	"strconv"
	"sync"
	"time"
)

type Service struct {
	log       *zerolog.Logger
	shift     int
	precision int64

	mu      sync.RWMutex
	storage map[domain.Ticker]map[int64]float64
}

func New(log *zerolog.Logger, shift int, precision int64) *Service {
	return &Service{
		storage:   make(map[domain.Ticker]map[int64]float64),
		shift:     shift,
		log:       log,
		precision: precision,
	}
}

func (s *Service) Put(ctx context.Context, ti domain.TickerPrice) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.log.Debug().Interface("price", ti).Msg("storage")

	select {
	case <-ctx.Done():
		return nil
	default:
		if _, ok := s.storage[ti.Ticker]; !ok {
			s.storage[ti.Ticker] = make(map[int64]float64)
		}

		tm := (ti.Time.Unix()/s.precision + 1) * s.precision

		stored := s.storage[ti.Ticker][tm]

		income, err := strconv.ParseFloat(ti.Price, 64)
		if err != nil {
			s.log.Err(err).Msg("try to parse val")
		} else {
			if income == 0 || s.shift < 2 {
				s.storage[ti.Ticker][tm] = stored
			} else {
				res := stored/float64(s.shift)*float64(s.shift-1) + income/float64(s.shift)

				s.log.Debug().Interface("value", res).Msg("storage")

				s.storage[ti.Ticker][tm] = res
			}
		}

		s.log.Debug().Interface("map", s.storage).Msg("storage")
	}

	return nil
}

func (s *Service) Get(ti domain.Ticker) (int64, float64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tmMarker := time.Now().Unix() / s.precision * s.precision

	tm, ok := s.storage[ti]
	if !ok {
		return 0, 0
	}

	s.log.Debug().Int64("tmn", tmMarker).Msg("time marker")

	return tmMarker, tm[tmMarker]
}

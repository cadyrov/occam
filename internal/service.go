package internal

import (
	"context"
	"github.com/cadyrov/occam/domain"
	"github.com/rs/zerolog"
	"sync"
	"time"
)

type Service struct {
	log *zerolog.Logger

	storage Keeper

	subscribers domain.SubscriberList
}

func New(log *zerolog.Logger, keeper Keeper) *Service {
	return &Service{
		log:     log,
		storage: keeper,
	}
}

func (s *Service) Run(ctx context.Context) {
	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		s.startSubscribers(ctx)

		wg.Done()
	}()

	wg.Add(1)

	go func() {
		s.result(ctx)

		wg.Done()
	}()

	wg.Wait()
}

func (s *Service) startSubscribers(ctx context.Context) {
	for i := range s.subscribers {

	}
}

func (s *Service) result(ctx context.Context) {

}

func (s *Service) Shutdown(ctx context.Context) error {
	s.log.Warn().Msg("server shutdown")

	time.Sleep(time.Second)

	return nil
}

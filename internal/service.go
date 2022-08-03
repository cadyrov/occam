package internal

import (
	"context"
	"fmt"
	"github.com/cadyrov/occam/domain"
	"github.com/rs/zerolog"
	"io"
	"sync"
	"time"
)

type Service struct {
	log *zerolog.Logger

	storage Keeper

	subscribers domain.SubscriberList

	precision int64
}

func New(log *zerolog.Logger, keeper Keeper, listOrigins []PriceStreamSubscriber, precision int64) *Service {
	srv := &Service{
		log:       log,
		storage:   keeper,
		precision: precision,
	}

	for i := range listOrigins {
		dmn := domain.Subscriber{}
		dmn.TTP, dmn.Err = listOrigins[i].SubscribePriceStream(domain.BTCUSDTicker)

		srv.subscribers = append(srv.subscribers, dmn)
	}

	return srv
}

func (s *Service) Run(ctx context.Context, output io.Writer) {
	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		s.startSubscribers(ctx)

		wg.Done()
	}()

	wg.Add(1)

	go func() {
		s.RunResult(ctx, output)

		wg.Done()
	}()

	wg.Wait()
}

func (s *Service) startSubscribers(ctx context.Context) {
	for i := range s.subscribers {
		sbs := s.subscribers[i]
		go func() {
			for {
				select {
				case <-ctx.Done():
					return

				case tk := <-sbs.TTP:
					if tk.Price != "" {
						s.storage.Put(ctx, tk)
					}
				case <-sbs.Err:
					return
				}
			}
		}()
	}
}

func (s *Service) RunResult(ctx context.Context, w io.Writer) {
	tm := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-tm.C:
			tm := time.Now()
			if tm.Unix()%s.precision != 0 {
				continue
			}

			timeMarker, val := s.storage.Get(domain.BTCUSDTicker)
			str := fmt.Sprintf("%d, %.2f \n", timeMarker, val)
			s.log.Debug().Str("value", str).Msg("try to write")

			if _, err := w.Write([]byte(str)); err != nil {
				s.log.Err(err).Msg("try to write")
			}

			s.storage.ClearOldest(domain.BTCUSDTicker, timeMarker)
		}
	}
}

package origin

import (
	"context"
	"fmt"
	"github.com/cadyrov/occam/domain"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"math/rand"
	"time"
)

type MockOrigin struct {
	log    *zerolog.Logger
	price  chan domain.TickerPrice
	err    chan error
	ticker *time.Ticker
	closed bool
}

func New(log *zerolog.Logger) *MockOrigin {
	return &MockOrigin{
		log:    log,
		price:  make(chan domain.TickerPrice),
		err:    make(chan error),
		ticker: time.NewTicker(time.Second),
	}
}

var (
	precisionToClosed = 10
	ErrCloseChan      = errors.New("some err")
)

func (mo *MockOrigin) Start(ctx context.Context) {
	go func() {
		tk := time.NewTicker(time.Second)
		for range tk.C {
			select {
			case <-ctx.Done():
				if !mo.closed {
					close(mo.price)
				}

				close(mo.err)
			default:
				if mo.closed {
					continue
				}

				if rand.Intn(precisionToClosed) == 1 {
					close(mo.price)

					mo.closed = true

					mo.err <- ErrCloseChan

					mo.log.Info().Str("value", fmt.Sprint(mo.price)).Msg("channel closed")

					continue
				}

				mo.price <- domain.TickerPrice{
					Ticker: domain.BTCUSDTicker,
					Time:   time.Now(),
					Price:  fmt.Sprintf("%f", rand.Float32()),
				}
			}
		}
	}()
}

func (mo *MockOrigin) SubscribePriceStream(domain.Ticker) (chan domain.TickerPrice, chan error) {
	return mo.price, mo.err
}

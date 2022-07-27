package origin

import (
	"github.com/cadyrov/occam/domain"
	"github.com/rs/zerolog"
	"time"
)

type MockOrigin struct {
	price  chan domain.TickerPrice
	err    chan error
	ticker *time.Ticker
}

func New(log zerolog.Logger) *MockOrigin {
	return &MockOrigin{
		price:  make(chan domain.TickerPrice),
		err:    make(chan error),
		ticker: time.NewTicker(time.Second),
	}
}

func (mo *MockOrigin) start() {

}

func (mo *MockOrigin) SubscribePriceStream(domain.Ticker) (chan domain.TickerPrice, chan error) {
	return mo.price, mo.err
}

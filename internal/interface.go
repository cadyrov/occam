package internal

import (
	"context"
	"github.com/cadyrov/occam/domain"
)

type PriceStreamSubscriber interface {
	SubscribePriceStream(domain.Ticker) (chan domain.TickerPrice, chan error)
}

type Keeper interface {
	Put(context.Context, domain.TickerPrice) error
	Get(ticker domain.Ticker) (int64, float64)
}

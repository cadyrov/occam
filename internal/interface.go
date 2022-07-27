package internal

import (
	"context"
	"github.com/cadyrov/occam/domain"
)

type PriceStreamSubscriber interface {
	SubscribePriceStream(domain.Ticker) (chan domain.TickerPrice, chan error)
}

type Keeper interface {
	Put(context.Context, domain.Ticker, float64) error
	Get(domain.Ticker) float64
}

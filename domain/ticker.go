package domain

import "time"

type Ticker string

const (
	BTCUSDTicker Ticker = "BTC_USD"
)

type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	Price  string // decimal value. example: "0", "10", "12.2", "13.2345122"
}

type Subscriber struct {
	TTP    chan TickerPrice
	Err    chan error
	Closed bool
}

type SubscriberList []Subscriber

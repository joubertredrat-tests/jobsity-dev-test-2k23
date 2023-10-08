package domain

import "context"

type StockQuote interface {
	Get(ctx context.Context, stock Stock) (Stock, error)
}

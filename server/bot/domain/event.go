package domain

import "context"

type StockEvent interface {
	StockRequested(ctx context.Context, stock Stock) error
}

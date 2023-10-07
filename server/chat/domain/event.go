package domain

import "context"

type MessageEvent interface {
	StockCommandReceived(ctx context.Context, message Message) error
}

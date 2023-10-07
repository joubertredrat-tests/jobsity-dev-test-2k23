package domain

import "context"

type MessageEvent interface {
	Created(ctx context.Context, message Message) (Message, error)
}

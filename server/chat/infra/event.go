package infra

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
)

type MessageEventRedis struct {
}

func NewMessageEventRedis() MessageEventRedis {
	return MessageEventRedis{}
}

func (e MessageEventRedis) Created(ctx context.Context, message domain.Message) (domain.Message, error) {
	return domain.Message{}, nil
}

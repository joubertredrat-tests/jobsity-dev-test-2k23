package infra

import (
	"context"
	"encoding/json"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type StockCommandReceived struct {
	Code       string `json:"code"`
	RawCommand string `json:"rawCommand"`
}

type MessageEventRedis struct {
	redisClient             *redis.Client
	logger                  *logrus.Logger
	messageCreatedTopicName string
}

func NewMessageEventRedis(c *redis.Client, l *logrus.Logger, mct string) domain.MessageEvent {
	return MessageEventRedis{
		redisClient:             c,
		logger:                  l,
		messageCreatedTopicName: mct,
	}
}

func (e MessageEventRedis) StockCommandReceived(ctx context.Context, message domain.Message) error {
	payload, err := json.Marshal(StockCommandReceived{
		Code:       message.StockCode(),
		RawCommand: message.Text,
	})
	if err != nil {
		e.logger.Error(err)
		return err
	}

	if err := e.redisClient.Publish(context.Background(), e.messageCreatedTopicName, payload).Err(); err != nil {
		e.logger.Error(err)
		return err
	}

	return nil
}

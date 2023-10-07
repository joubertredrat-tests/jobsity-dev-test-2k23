package infra

import (
	"context"
	"encoding/json"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"joubertredrat-tests/jobsity-dev-test-2k23/pkg"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type MessageCreated struct {
	ID          string  `json:"id"`
	UserName    string  `json:"userName"`
	UserEmail   string  `json:"userEmail"`
	MessageText string  `json:"messageText"`
	Datetime    *string `json:"datetime"`
}

type MessageEventRedis struct {
	redisClient             *redis.Client
	logger                  *logrus.Logger
	messageCreatedTopicName string
}

func NewMessageEventRedis(c *redis.Client, l *logrus.Logger, mct string) MessageEventRedis {
	return MessageEventRedis{
		redisClient:             c,
		logger:                  l,
		messageCreatedTopicName: mct,
	}
}

func (e MessageEventRedis) Created(ctx context.Context, message domain.Message) (domain.Message, error) {
	payload, err := json.Marshal(MessageCreated{
		ID:          message.ID,
		UserName:    message.UserName,
		UserEmail:   message.UserEmail,
		MessageText: message.Text,
		Datetime:    pkg.DatetimeCanonical(&message.Datetime),
	})
	if err != nil {
		e.logger.Error(err)
		return domain.Message{}, err
	}

	if err := e.redisClient.Publish(context.Background(), e.messageCreatedTopicName, payload).Err(); err != nil {
		e.logger.Error(err)
		return domain.Message{}, err
	}

	return message, nil
}

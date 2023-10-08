package infra

import (
	"context"
	"encoding/json"
	"joubertredrat-tests/jobsity-dev-test-2k23/bot/domain"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type StockCommandRequested struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}

type StockEventRedis struct {
	redisClient *redis.Client
	logger      *logrus.Logger
	topicName   string
}

func NewStockEventRedis(c *redis.Client, l *logrus.Logger, t string) domain.StockEvent {
	return StockEventRedis{
		redisClient: c,
		logger:      l,
		topicName:   t,
	}
}

func (e StockEventRedis) StockRequested(ctx context.Context, stock domain.Stock) error {
	payload, err := json.Marshal(StockCommandRequested{
		Code:  stock.Code,
		Value: stock.Value,
	})
	if err != nil {
		e.logger.Error(err)
		return err
	}

	if err := e.redisClient.Publish(context.Background(), e.topicName, payload).Err(); err != nil {
		e.logger.Error(err)
		return err
	}

	return nil
}

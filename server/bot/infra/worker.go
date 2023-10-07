package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"joubertredrat-tests/jobsity-dev-test-2k23/bot/application"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type StockCommandReceived struct {
	Code       string `json:"code"`
	RawCommand string `json:"rawCommand"`
}

type BotWorker struct {
	usecase            application.UsecaseGetStockValue
	redisClient        *redis.Client
	logger             *logrus.Logger
	subscribeTopicName string
}

func NewBotWorker(
	u application.UsecaseGetStockValue,
	c *redis.Client,
	l *logrus.Logger,
	s string,
) BotWorker {
	return BotWorker{
		usecase:            u,
		redisClient:        c,
		logger:             l,
		subscribeTopicName: s,
	}
}

func (w BotWorker) Run(ctx context.Context) {
	subscriber := w.redisClient.Subscribe(ctx, w.subscribeTopicName)

	fmt.Println("Worker running")

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			w.logger.Error(err)
		}

		var command StockCommandReceived
		if err := json.Unmarshal([]byte(msg.Payload), &command); err != nil {
			w.logger.Error(err)
		}

		stock, err := w.usecase.Execute(ctx, application.UsecaseGetStockValueInput{
			Code: command.Code,
		})
		if err != nil {
			w.logger.Error(err)
		}

		w.logger.Info(fmt.Sprintf("Stock requested: Code [ %s ] Value [ %s ]", stock.Code, stock.Value))
	}
}

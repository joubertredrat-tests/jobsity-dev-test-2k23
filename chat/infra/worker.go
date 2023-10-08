package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/application"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const STOCK_NO_QUOTE = "N/D"

type StockCommandRequested struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}

type StockWorker struct {
	usecase            application.UsecaseMessageCreate
	redisClient        *redis.Client
	logger             *logrus.Logger
	subscribeTopicName string
}

func NewStockWorker(
	u application.UsecaseMessageCreate,
	c *redis.Client,
	l *logrus.Logger,
	s string,
) StockWorker {
	return StockWorker{
		usecase:            u,
		redisClient:        c,
		logger:             l,
		subscribeTopicName: s,
	}
}

func (w StockWorker) Run(ctx context.Context) {
	subscriber := w.redisClient.Subscribe(ctx, w.subscribeTopicName)

	fmt.Println("Worker running")

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			w.logger.Error(err)
		}

		var command StockCommandRequested
		if err := json.Unmarshal([]byte(msg.Payload), &command); err != nil {
			w.logger.Error(err)
		}

		text := fmt.Sprintf("%s quote is %s per share", command.Code, command.Value)
		if command.Value == STOCK_NO_QUOTE {
			text = fmt.Sprintf("%s has no quote, check if code is correct", command.Code)
		}

		message, err := w.usecase.Execute(ctx, application.UsecaseMessageCreateInput{
			UserName:    "Stock Bot",
			UserEmail:   "stock.bot@jobsity.com",
			MessageText: text,
		})
		if err != nil {
			w.logger.Error(err)
		}

		w.logger.Info(fmt.Sprintf(
			"Stock message ID %s created by %s <%s> %s",
			message.ID,
			message.UserName,
			message.UserEmail,
			message.Text,
		))
	}
}

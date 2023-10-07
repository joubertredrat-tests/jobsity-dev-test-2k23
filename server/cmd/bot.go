package cmd

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/bot/application"
	"joubertredrat-tests/jobsity-dev-test-2k23/bot/infra"

	"github.com/urfave/cli/v2"
)

func getBotCommand() *cli.Command {
	return &cli.Command{
		Name:    "bot",
		Aliases: []string{},
		Usage:   "Run bot worker to get stock values",
		Action: func(c *cli.Context) error {
			ctx := context.TODO()

			config, err := infra.NewConfig()
			if err != nil {
				return err
			}

			logger := infra.Logger()

			redis := infra.RedisClient(config.RedisQueueHost, config.RedisQueuePort)

			stockService := infra.NewStockQuoteStooq()
			stockEvent := infra.NewStockEventRedis(redis, logger, config.RedisQueueStockRequestedTopicName)

			usecaseGetStockValue := application.NewUsecaseGetStockValue(stockService, stockEvent)

			worker := infra.NewBotWorker(usecaseGetStockValue, redis, logger, config.RedisQueueStockCommandReceivedTopicName)
			worker.Run(ctx)

			return nil
		},
	}
}

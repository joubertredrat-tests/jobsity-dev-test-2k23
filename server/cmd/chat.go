package cmd

import (
	"context"
	"fmt"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/application"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/infra"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func getChatApiCommand() *cli.Command {
	return &cli.Command{
		Name:    "chat",
		Aliases: []string{},
		Usage:   "Open chat HTTP api to listen",
		Action: func(c *cli.Context) error {
			config, err := infra.NewConfig()
			if err != nil {
				return err
			}

			r := gin.Default()
			if err := r.SetTrustedProxies(nil); err != nil {
				return err
			}
			corsConfig := cors.DefaultConfig()
			corsConfig.AllowOrigins = []string{"http://0.0.0.0:9090", "http://127.0.0.1:9090"}
			corsConfig.AllowMethods = []string{"GET", "POST", "OPTIONS"}
			corsConfig.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
			r.Use(cors.New(corsConfig))

			logger := infra.Logger()
			dsn := infra.MongoDSN(
				config.MongoHost,
				config.MongoPort,
				config.MongoUser,
				config.MongoPassword,
			)
			mongo, err := infra.MongoConnection(dsn)
			if err != nil {
				return err
			}
			redis := infra.RedisClient(config.RedisQueueHost, config.RedisQueuePort)

			tokenService := infra.NewTokenServiceJWT(
				logger,
				config.JwtSecretKey,
				config.JwtTokenExpirationHours,
			)

			userRepository := infra.NewUserRepository(mongo, config.MongoDatabase, logger)
			messageRepository := infra.NewMessageRepository(mongo, config.MongoDatabase, logger)
			messageEvent := infra.NewMessageEventRedis(redis, logger, config.RedisQueueStockCommandReceivedTopicName)

			usecaseUserRegister := application.NewUsecaseUserRegister(userRepository)
			usecaseUserLogin := application.NewUsecaseUserLogin(userRepository, tokenService)
			usecaseMessageCreate := application.NewUsecaseMessageCreate(messageRepository, messageEvent)
			usecaseMessagesList := application.NewUsecaseMessageList(messageRepository)

			apiBaseController := infra.NewApiBaseController()
			userController := infra.NewUserController()
			messagesController := infra.NewMessagesController()

			r.NoRoute(apiBaseController.HandleNotFound)

			ra := r.Group("/api")
			infra.RegisterCustomValidator()
			{
				ra.GET("/status", apiBaseController.HandleStatus)
				ra.POST(
					"/register",
					infra.JSONBodyMiddleware(),
					userController.HandleCreate(usecaseUserRegister),
				)
				ra.POST(
					"/login",
					infra.JSONBodyMiddleware(),
					userController.HandleLogin(usecaseUserLogin),
				)
				ra.POST(
					"/messages",
					infra.JwtCheckMiddleware(tokenService),
					infra.JSONBodyMiddleware(),
					messagesController.HandleCreate(usecaseMessageCreate),
				)
				ra.GET(
					"/messages",
					infra.JwtCheckMiddleware(tokenService),
					messagesController.HandleList(usecaseMessagesList),
				)
			}

			worker := infra.NewStockWorker(
				usecaseMessageCreate,
				redis,
				logger,
				config.RedisQueueStockRequestedTopicName,
			)

			go worker.Run(context.TODO())

			return r.Run(fmt.Sprintf("%s:%s", config.ApiHost, config.ApiPort))
		},
	}
}

package cmd

import (
	"fmt"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/application"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/infra"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func getApiCommand() *cli.Command {
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

			userRepository := infra.NewUserRepository(mongo, logger)

			usecaseUserRegister := application.NewUsecaseUserRegister(userRepository)

			apiBaseController := infra.NewApiBaseController()
			userController := infra.NewUserController()

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
			}

			return r.Run(fmt.Sprintf("%s:%s", config.ApiHost, config.ApiPort))
		},
	}
}

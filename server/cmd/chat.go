package cmd

import (
	"fmt"
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

			apiBaseController := infra.NewApiBaseController()

			r.NoRoute(apiBaseController.HandleNotFound)

			ra := r.Group("/api")
			{
				ra.GET("/status", apiBaseController.HandleStatus)
			}

			return r.Run(fmt.Sprintf("%s:%s", config.ApiHost, config.ApiPort))
		},
	}
}

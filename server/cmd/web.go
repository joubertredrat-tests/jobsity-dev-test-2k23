package cmd

import (
	"fmt"
	"joubertredrat-tests/jobsity-dev-test-2k23/web"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func getWebCommand() *cli.Command {
	return &cli.Command{
		Name:    "web",
		Aliases: []string{},
		Usage:   "Start web interface",
		Action: func(c *cli.Context) error {
			config, err := web.NewConfig()
			if err != nil {
				return err
			}

			r := gin.Default()
			if err := r.SetTrustedProxies(nil); err != nil {
				return err
			}

			r.LoadHTMLGlob("static/*.html")
			r.StaticFile("/favicon.ico", "static/favicon.ico")
			r.StaticFile("/login.js", "static/login.js")
			r.GET("/login", func(c *gin.Context) {
				chatToken, _ := c.Cookie("chatToken")
				if chatToken != "" {
					c.Redirect(http.StatusTemporaryRedirect, "/chat")
					return
				}

				c.HTML(http.StatusOK, "login.html", nil)
			})
			r.GET("/register", func(c *gin.Context) {
				chatToken, _ := c.Cookie("chatToken")
				if chatToken != "" {
					c.Redirect(http.StatusTemporaryRedirect, "/chat")
					return
				}

				c.HTML(http.StatusOK, "register.html", nil)
			})
			r.GET("/chat", func(c *gin.Context) {
				_, err := c.Cookie("chatToken")
				if err != nil {
					c.Redirect(http.StatusTemporaryRedirect, "/login")
					return
				}

				c.HTML(http.StatusOK, "chat.html", nil)
			})

			return r.Run(fmt.Sprintf("%s:%s", config.WebHost, config.WebPort))
		},
	}
}

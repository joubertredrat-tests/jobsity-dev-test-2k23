package infra

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func RedisClient(host, port string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
}

func Logger() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	return log
}

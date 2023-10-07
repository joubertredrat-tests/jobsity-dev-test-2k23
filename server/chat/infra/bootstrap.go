package infra

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDSN(host, port, user, password string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
}

func MongoConnection(dsn string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return mongo.Connect(ctx, options.Client().ApplyURI(dsn))
}

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

package infra

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryMongo struct {
	mongoClient *mongo.Client
	logger      *logrus.Logger
}

func NewUserRepository(c *mongo.Client, l *logrus.Logger) domain.UserRepository {
	return UserRepositoryMongo{
		mongoClient: c,
		logger:      l,
	}
}

func (r UserRepositoryMongo) Persist(ctx context.Context, user domain.User) (domain.User, error) {
	return domain.User{
		ID:       "ID",
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (r UserRepositoryMongo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return domain.User{}, nil
}

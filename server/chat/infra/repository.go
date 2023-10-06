package infra

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	COLLECTION_USERS = "users"
)

type UserRepositoryMongo struct {
	mongoClient *mongo.Client
	database    string
	logger      *logrus.Logger
}

type UserMongo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

func NewUserRepository(c *mongo.Client, d string, l *logrus.Logger) domain.UserRepository {
	return UserRepositoryMongo{
		mongoClient: c,
		database:    d,
		logger:      l,
	}
}

func (r UserRepositoryMongo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	collection := r.collection()
	hashPassword, err := hashPassword([]byte(user.Password))
	if err != nil {
		r.logger.Error(err)
		return domain.User{}, err
	}

	userMongo := UserMongo{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashPassword,
	}

	result, err := collection.InsertOne(ctx, userMongo)
	if err != nil {
		r.logger.Error(err)
		return domain.User{}, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	user.Password = hashPassword
	return user, nil
}

func (r UserRepositoryMongo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	collection := r.collection()

	var user UserMongo
	err := collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		r.logger.Error(err)
		return domain.User{}, nil
	}
	if err != nil {
		r.logger.Error(err)
		return domain.User{}, err
	}

	return domain.User{
		ID:       user.ID.Hex(),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (r UserRepositoryMongo) collection() *mongo.Collection {
	return r.mongoClient.Database(r.database).Collection(COLLECTION_USERS)
}

func hashPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

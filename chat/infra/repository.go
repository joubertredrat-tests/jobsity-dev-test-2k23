package infra

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	COLLECTION_USERS    = "users"
	COLLECTION_MESSAGES = "messages"
)

type UserMongo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type UserRepositoryMongo struct {
	mongoClient *mongo.Client
	database    string
	logger      *logrus.Logger
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

	var userMongo UserMongo
	err := collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&userMongo)
	if err == mongo.ErrNoDocuments {
		r.logger.Error(err)
		return domain.User{}, nil
	}
	if err != nil {
		r.logger.Error(err)
		return domain.User{}, err
	}

	return domain.User{
		ID:       userMongo.ID.Hex(),
		Name:     userMongo.Name,
		Email:    userMongo.Email,
		Password: userMongo.Password,
	}, nil
}

func (r UserRepositoryMongo) GetAuthenticated(ctx context.Context, user domain.User) (domain.User, error) {
	collection := r.collection()

	var userMongo UserMongo
	err := collection.FindOne(ctx, bson.D{{"email", user.Email}}).Decode(&userMongo)
	if err == mongo.ErrNoDocuments {
		r.logger.Error(err)
		return domain.User{}, domain.NewErrUserNotFoundByEmail(user.Email)
	}
	if err != nil {
		r.logger.Error(err)
		return domain.User{}, err
	}

	errPass := bcrypt.CompareHashAndPassword([]byte(userMongo.Password), []byte(user.Password))
	if errPass != nil {
		r.logger.Error(err)
		return domain.User{}, domain.NewErrUserNotAuthenticated(user.Email)
	}

	return domain.User{
		ID:       userMongo.ID.Hex(),
		Name:     userMongo.Name,
		Email:    userMongo.Email,
		Password: userMongo.Password,
	}, nil
}

func (r UserRepositoryMongo) collection() *mongo.Collection {
	return r.mongoClient.Database(r.database).Collection(COLLECTION_USERS)
}

type MessageMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AppID     string             `bson:"appId"`
	UserName  string             `bson:"name"`
	UserEmail string             `bson:"email"`
	Text      string             `bson:"text"`
	Datetime  time.Time          `bson:"datetime"`
}

type MessageRepositoryMongo struct {
	mongoClient *mongo.Client
	database    string
	logger      *logrus.Logger
}

func NewMessageRepository(c *mongo.Client, d string, l *logrus.Logger) domain.MessageRepository {
	return MessageRepositoryMongo{
		mongoClient: c,
		database:    d,
		logger:      l,
	}
}

func (r MessageRepositoryMongo) Create(ctx context.Context, message domain.Message) (domain.Message, error) {
	collection := r.collection()

	messageMongo := MessageMongo{
		AppID:     ulid.Make().String(),
		UserName:  message.UserName,
		UserEmail: message.UserEmail,
		Text:      message.Text,
		Datetime:  time.Now(),
	}

	_, err := collection.InsertOne(ctx, messageMongo)
	if err != nil {
		r.logger.Error(err)
		return domain.Message{}, err
	}

	message.ID = messageMongo.AppID
	message.Datetime = messageMongo.Datetime
	return message, nil
}

func (r MessageRepositoryMongo) List(ctx context.Context, pagination domain.Pagination) ([]domain.Message, error) {
	collection := r.collection()
	pageOptions := skipLimit(pagination)

	sort := bson.D{}
	sort = append(sort, bson.E{"datetime", -1})
	pageOptions.SetSort(sort)

	cursor, err := collection.Find(ctx, bson.D{{}}, pageOptions)
	if err != nil {
		r.logger.Error(err)
		return []domain.Message{}, err
	}
	defer cursor.Close(ctx)

	list := []domain.Message{}

	for cursor.Next(ctx) {
		var messageMongo MessageMongo
		if err := cursor.Decode(&messageMongo); err != nil {
			r.logger.Error(err)
			return []domain.Message{}, err
		}
		list = append(list, domain.Message{
			ID:        messageMongo.AppID,
			UserName:  messageMongo.UserName,
			UserEmail: messageMongo.UserEmail,
			Text:      messageMongo.Text,
			Datetime:  messageMongo.Datetime,
		})
	}

	if err := cursor.Err(); err != nil {
		r.logger.Error(err)
		return []domain.Message{}, err
	}

	return list, nil
}

func skipLimit(pagination domain.Pagination) *options.FindOptions {
	pageOptions := options.Find()
	pageOptions.SetSkip(int64(pagination.Page - 1))
	pageOptions.SetLimit(int64(pagination.ItemsPerPage))

	return pageOptions
}

func (r MessageRepositoryMongo) collection() *mongo.Collection {
	return r.mongoClient.Database(r.database).Collection(COLLECTION_MESSAGES)
}

func hashPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

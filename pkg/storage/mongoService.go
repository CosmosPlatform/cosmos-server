package storage

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/storage/obj"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoService struct {
	mongo              *mongo.Client
	databaseName       string
	userCollectionName string
	logger             log.Logger
}

func NewMongoService(config config.StorageConfig, logger log.Logger) (*MongoService, error) {
	uri := fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)

	clientOptions := options.Client().ApplyURI(uri)
	mongoClient, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &MongoService{
		mongo:              mongoClient,
		databaseName:       "cosmos",
		userCollectionName: "users",
		logger:             logger,
	}, nil
}

func (s *MongoService) GetUserCollection() *mongo.Collection {
	return s.mongo.Database(s.databaseName).Collection(s.userCollectionName)
}

func (s *MongoService) InsertUser(ctx context.Context, user *obj.User) error {
	collection := s.GetUserCollection()
	_, err := insertOne[obj.User](ctx, collection, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (s *MongoService) GetUserWithEmail(ctx context.Context, email string) (*obj.User, error) {
	collection := s.GetUserCollection()

	filter := NewBsonBuilder().
		Eq("email", email).
		Build()

	user, err := findOne[obj.User](ctx, collection, filter)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *MongoService) GetUserWithRole(ctx context.Context, role string) (*obj.User, error) {
	collection := s.GetUserCollection()

	filter := NewBsonBuilder().
		Eq("role", role).
		Build()

	user, err := findOne[obj.User](ctx, collection, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user with role %s: %w", role, err)
	}
	return user, nil
}

package storage

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/storage/obj"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoService struct {
	mongo              *mongo.Client
	databaseName       string
	userCollectionName string
}

func NewMongoService(config config.StorageConfig) (*MongoService, error) {
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
	}, nil
}

func (s *MongoService) GetUserCollection() *mongo.Collection {
	return s.mongo.Database(s.databaseName).Collection(s.userCollectionName)
}

func (s *MongoService) InsertUser(ctx context.Context, user *obj.User) (string, error) {
	collection := s.GetUserCollection()
	result, err := insertOne[obj.User](ctx, collection, user)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %w", err)
	}
	return result.InsertedID.(string), nil
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

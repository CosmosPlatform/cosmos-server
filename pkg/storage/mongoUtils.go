package storage

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BsonBuilder struct {
	filter bson.M
}

// NewBsonBuilder creates a new builder instance.
func NewBsonBuilder() *BsonBuilder {
	return &BsonBuilder{filter: bson.M{}}
}

// Eq adds an equality condition.
func (b *BsonBuilder) Eq(field string, value interface{}) *BsonBuilder {
	b.filter[field] = value
	return b
}

// Gt adds a greater-than condition.
func (b *BsonBuilder) Gt(field string, value interface{}) *BsonBuilder {
	b.ensureFieldIsMap(field)
	b.filter[field].(bson.M)["$gt"] = value
	return b
}

// Lt adds a less-than condition.
func (b *BsonBuilder) Lt(field string, value interface{}) *BsonBuilder {
	b.ensureFieldIsMap(field)
	b.filter[field].(bson.M)["$lt"] = value
	return b
}

// In adds an $in condition.
func (b *BsonBuilder) In(field string, values ...interface{}) *BsonBuilder {
	b.filter[field] = bson.M{"$in": values}
	return b
}

// And adds an $and block.
func (b *BsonBuilder) And(conditions ...bson.M) *BsonBuilder {
	b.filter["$and"] = conditions
	return b
}

// Or adds an $or block.
func (b *BsonBuilder) Or(conditions ...bson.M) *BsonBuilder {
	b.filter["$or"] = conditions
	return b
}

// Build returns the final bson.M filter.
func (b *BsonBuilder) Build() bson.M {
	return b.filter
}

// ensureFieldIsMap ensures that the field is a bson.M so we can chain operators.
func (b *BsonBuilder) ensureFieldIsMap(field string) {
	if _, ok := b.filter[field]; !ok {
		b.filter[field] = bson.M{}
	}
}

// findOne retrieves a single document from the given collection and decodes it into T.
func findOne[T any](ctx context.Context, collection *mongo.Collection, filter interface{}) (*T, error) {
	var result T
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return &result, nil
}

// FindMany retrieves all matching documents and decodes them into a slice of T.
func findMany[T any](ctx context.Context, collection *mongo.Collection, filter interface{}) ([]T, error) {
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Printf("Error closing cursor: %v\n", err)
		}
	}()

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// InsertOne inserts a single document of type T into the given collection.
func insertOne[T any](ctx context.Context, collection *mongo.Collection, doc *T) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}
	return result, nil
}

// InsertMany inserts multiple documents of type T into the given collection.
func insertMany[T any](ctx context.Context, collection *mongo.Collection, docs []T) (*mongo.InsertManyResult, error) {
	// Convert []T to []interface{}
	values := make([]interface{}, len(docs))
	for i := range docs {
		values[i] = docs[i]
	}

	result, err := collection.InsertMany(ctx, values)
	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %v", err)
	}
	return result, nil
}

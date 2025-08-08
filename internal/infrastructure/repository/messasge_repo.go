package repository

import (
	"context"

	"github.com/noxhalley/funken/internal/infrastructure/log"
	"github.com/noxhalley/funken/internal/infrastructure/mongodb"
	"github.com/noxhalley/funken/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MessageRepository interface {
	CountByConditions(
		ctx context.Context,
		filter interface{},
		opts *options.CountOptionsBuilder,
	) (int64, error)

	FindOneByConditions(
		ctx context.Context,
		filter interface{},
		opts *options.FindOneOptionsBuilder,
	) (*model.Message, error)

	FindByConditions(
		ctx context.Context,
		filter interface{},
		opts *options.FindOptionsBuilder,
	) ([]model.Message, error)

	Create(
		ctx context.Context,
		msg model.Message,
	) error

	UpdateByID(
		ctx context.Context,
		ID string,
		operation interface{},
	) (*model.Message, error)

	DeleteByID(
		ctx context.Context,
		ID string,
	) error
}

type messageRepo struct {
	logger *log.Logger
	coll   *mongo.Collection
}

func NewMessageRepository(db *mongodb.MongoDB) MessageRepository {
	coll := db.Client.
		Database(db.DBName).
		Collection(model.MessageCollectionName)

	return &messageRepo{
		logger: log.With("repository", "message_repository"),
		coll:   coll,
	}
}

// CountByConditions implements MessageRepository.
func (m *messageRepo) CountByConditions(
	ctx context.Context,
	filter interface{},
	opts *options.CountOptionsBuilder,
) (int64, error) {
	return m.coll.CountDocuments(ctx, filter, opts)
}

// FindOneByConditions implements MessageRepository.
func (m *messageRepo) FindOneByConditions(
	ctx context.Context,
	filter interface{},
	opts *options.FindOneOptionsBuilder,
) (*model.Message, error) {
	msg := model.Message{}
	if err := m.coll.FindOne(ctx, filter, opts).Decode(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// FindByConditions implements MessageRepository.
func (m *messageRepo) FindByConditions(
	ctx context.Context,
	filter interface{},
	opts *options.FindOptionsBuilder,
) ([]model.Message, error) {
	cursor, err := m.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []model.Message
	err = cursor.All(ctx, &messages)
	return messages, err
}

// Create implements MessageRepository.
func (m *messageRepo) Create(
	ctx context.Context,
	msg model.Message,
) error {
	_, err := m.coll.InsertOne(ctx, msg)
	return err
}

// UpdateByID implements MessageRepository.
func (m *messageRepo) UpdateByID(
	ctx context.Context,
	ID string,
	operation interface{},
) (*model.Message, error) {
	filter := bson.M{"id": ID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	updatedDoc := model.Message{}
	if err := m.coll.FindOneAndUpdate(ctx, filter, operation, opts).Decode(&updatedDoc); err != nil {
		return nil, err
	}
	return &updatedDoc, nil
}

// DeleteByID implements MessageRepository.
func (m *messageRepo) DeleteByID(ctx context.Context, ID string) error {
	filter := bson.M{"id": ID}
	_, err := m.coll.DeleteOne(ctx, filter)
	return err
}

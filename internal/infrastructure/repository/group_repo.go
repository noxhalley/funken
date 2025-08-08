package repository

import (
	"context"
	"errors"

	"github.com/noxhalley/funken/internal/infrastructure/log"
	"github.com/noxhalley/funken/internal/infrastructure/mongodb"
	"github.com/noxhalley/funken/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrInvalidDocumentID = errors.New("Document's ID is invalid")
)

type GroupRepository interface {
	FindOneByConditions(
		ctx context.Context,
		filter interface{},
		opts *options.FindOneOptionsBuilder,
	) (*model.Group, error)

	FindByConditions(
		ctx context.Context,
		filter interface{},
		opts *options.FindOptionsBuilder,
	) ([]model.Group, error)

	Create(
		ctx context.Context,
		group model.Group,
	) error

	UpdateByID(
		ctx context.Context,
		ID string,
		operation interface{},
	) (*model.Group, error)

	DeleteByID(
		ctx context.Context,
		ID string,
	) error

	CheckExist(ctx context.Context, ID string) (bool, error)
}

type groupRepo struct {
	logger *log.Logger
	coll   *mongo.Collection
}

func NewGroupRepository(db *mongodb.MongoDB) GroupRepository {
	coll := db.Client.
		Database(db.DBName).
		Collection(model.GroupCollectionName)

	return &groupRepo{
		logger: log.With("repository", "group_repository"),
		coll:   coll,
	}
}

// FindOneByConditions implements GroupRepository.
func (g *groupRepo) FindOneByConditions(
	ctx context.Context,
	filter interface{},
	opts *options.FindOneOptionsBuilder,
) (*model.Group, error) {
	group := model.Group{}
	if err := g.coll.FindOne(ctx, filter, opts).Decode(&group); err != nil {
		return nil, err
	}
	return &group, nil
}

// FindByConditions implements GroupRepository.
func (g *groupRepo) FindByConditions(
	ctx context.Context,
	filter interface{},
	opts *options.FindOptionsBuilder,
) ([]model.Group, error) {
	cursor, err := g.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []model.Group
	err = cursor.All(ctx, &groups)
	return groups, err
}

// Create implements GroupRepository.
func (g *groupRepo) Create(
	ctx context.Context,
	group model.Group,
) error {
	_, err := g.coll.InsertOne(ctx, group)
	return err
}

// UpdateByID implements GroupRepository.
func (g *groupRepo) UpdateByID(
	ctx context.Context,
	ID string,
	operation interface{},
) (*model.Group, error) {
	filter := bson.M{"id": ID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	updatedDoc := model.Group{}
	if err := g.coll.FindOneAndUpdate(ctx, filter, operation, opts).Decode(&updatedDoc); err != nil {
		return nil, err
	}
	return &updatedDoc, nil
}

// DeleteByID implements GroupRepository.
func (g *groupRepo) DeleteByID(ctx context.Context, ID string) error {
	filter := bson.M{"id": ID}
	_, err := g.coll.DeleteOne(ctx, filter)
	return err
}

// CheckExist implements GroupRepository.
func (g *groupRepo) CheckExist(ctx context.Context, ID string) (bool, error) {
	filter := bson.M{"id": ID}
	opts := options.
		FindOne().
		SetProjection(bson.M{"id": 1})

	err := g.coll.FindOne(ctx, filter, opts).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

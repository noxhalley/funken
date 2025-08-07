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
	FindOneByConditions(ctx context.Context, filter interface{}, opts *options.FindOneOptionsBuilder) (*model.Group, error)
	Create(ctx context.Context, group model.Group) (*model.Group, error)
	UpdateByID(ctx context.Context, ID string, operation interface{}) (*model.Group, error)
	CheckExist(ctx context.Context, ID string) (bool, error)
}

type groupRepo struct {
	logger *log.Logger
	db     *mongodb.MongoDB
}

func NewGroupRepository(db *mongodb.MongoDB) GroupRepository {
	return &groupRepo{
		logger: log.With("repository", "group_repository"),
		db:     db,
	}
}

// FindOneByConditions implements GroupRepository.
func (g *groupRepo) FindOneByConditions(
	ctx context.Context,
	filter interface{},
	opts *options.FindOneOptionsBuilder,
) (*model.Group, error) {
	group := model.Group{}
	coll := g.db.Client.
		Database(g.db.DBName).
		Collection(model.GroupCollectionName)

	if err := coll.FindOne(ctx, filter, opts).Decode(&group); err != nil {
		return nil, err
	}

	return &group, nil
}

// Create implements GroupRepository.
func (g *groupRepo) Create(ctx context.Context, group model.Group) (*model.Group, error) {
	panic("unimplemented")
}

// UpdateByID implements GroupRepository.
func (g *groupRepo) UpdateByID(
	ctx context.Context,
	ID string,
	operation interface{},
) (*model.Group, error) {
	filter := bson.M{"id": ID}
	coll := g.db.Client.
		Database(g.db.DBName).
		Collection(model.GroupCollectionName)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	updatedDoc := model.Group{}
	if err := coll.FindOneAndUpdate(ctx, filter, operation, opts).Decode(&updatedDoc); err != nil {
		return nil, err
	}
	return &updatedDoc, nil
}

// CheckExist implements GroupRepository.
func (g *groupRepo) CheckExist(ctx context.Context, ID string) (bool, error) {
	coll := g.db.Client.
		Database(g.db.DBName).
		Collection(model.GroupCollectionName)

	filter := bson.M{"id": ID}
	opts := options.
		FindOne().
		SetProjection(bson.M{"id": 1})

	err := coll.FindOne(ctx, filter, opts).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

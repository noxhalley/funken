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

type GroupNGFilterRepository interface {
	FindOneByConditions(
		ctx context.Context,
		filter interface{},
		opts *options.FindOneOptionsBuilder,
	) (*model.GroupNGFilter, error)

	FindByConditions(
		ctx context.Context,
		filter interface{},
		opts *options.FindOptionsBuilder,
	) ([]model.GroupNGFilter, error)

	Create(
		ctx context.Context,
		ngFilter model.GroupNGFilter,
	) error

	CreateBatch(
		ctx context.Context,
		ngFilters []model.GroupNGFilter,
	) error

	UpdateByID(
		ctx context.Context,
		ID string,
		operation interface{},
	) (*model.GroupNGFilter, error)

	DeleteByID(
		ctx context.Context,
		ID string,
	) error

	DeleteByGroupIDs(
		ctx context.Context,
		groupIDs []string,
	) error
}

type groupNGFilterRepo struct {
	logger *log.Logger
	coll   *mongo.Collection
}

func NewGroupNGFilterRepository(db *mongodb.MongoDB) GroupNGFilterRepository {
	coll := db.Client.
		Database(db.DBName).
		Collection(model.GroupNGFilterCollectionName)

	return &groupNGFilterRepo{
		logger: log.With("repository", "group_ng_filter_repository"),
		coll:   coll,
	}
}

// FindOneByConditions implements GroupNGFilterRepository.
func (g *groupNGFilterRepo) FindOneByConditions(
	ctx context.Context,
	filter interface{},
	opts *options.FindOneOptionsBuilder,
) (*model.GroupNGFilter, error) {
	ngFilter := model.GroupNGFilter{}
	if err := g.coll.FindOne(ctx, filter, opts).Decode(&ngFilter); err != nil {
		return nil, err
	}
	return &ngFilter, nil
}

// FindByConditions implements GroupNGFilterRepository.
func (g *groupNGFilterRepo) FindByConditions(
	ctx context.Context,
	filter interface{},
	opts *options.FindOptionsBuilder,
) ([]model.GroupNGFilter, error) {
	cursor, err := g.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ngFilters []model.GroupNGFilter
	err = cursor.All(ctx, &ngFilters)
	return ngFilters, err
}

// Create implements GroupNGFilterRepository.
func (g *groupNGFilterRepo) Create(
	ctx context.Context,
	ngFilter model.GroupNGFilter,
) error {
	_, err := g.coll.InsertOne(ctx, ngFilter)
	return err
}

// CreateBatch implements GroupNGFilterRepository.
func (g *groupNGFilterRepo) CreateBatch(
	ctx context.Context,
	ngFilters []model.GroupNGFilter,
) error {
	_, err := g.coll.InsertMany(ctx, ngFilters)
	return err
}

// UpdateByID implements GroupNGFilterRepository.
func (g *groupNGFilterRepo) UpdateByID(
	ctx context.Context,
	ID string,
	operation interface{},
) (*model.GroupNGFilter, error) {
	filter := bson.M{"id": ID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	updatedDoc := model.GroupNGFilter{}

	if err := g.coll.FindOneAndUpdate(ctx, filter, operation, opts).Decode(&updatedDoc); err != nil {
		return nil, err
	}
	return &updatedDoc, nil
}

// DeleteByID implements GroupNGFilterRepository.
func (g *groupNGFilterRepo) DeleteByID(
	ctx context.Context,
	ID string,
) error {
	filter := bson.M{"id": ID}
	_, err := g.coll.DeleteOne(ctx, filter)
	return err
}

// DeleteByGroupIDs implements GroupNGFilterRepository.
func (g *groupNGFilterRepo) DeleteByGroupIDs(
	ctx context.Context,
	groupIDs []string,
) error {
	filter := bson.M{
		"group_id": bson.M{
			"$in": groupIDs,
		},
	}

	_, err := g.coll.DeleteMany(ctx, filter)
	return err
}

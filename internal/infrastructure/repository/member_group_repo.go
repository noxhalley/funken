package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/noxhalley/funken/internal/infrastructure/log"
	"github.com/noxhalley/funken/internal/infrastructure/mongodb"
	"github.com/noxhalley/funken/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MemberGroupRepository interface {
	FindMemberIDsByGroupID(ctx context.Context, groupID string) ([]string, error)

	CountMembersByGroupID(ctx context.Context, groupID string) (int64, error)

	AddMembers(ctx context.Context, groupID string, memberIDs []string) error

	RemoveMembers(ctx context.Context, groupID string, memberIDs []string) error
}

type memberGroupRepo struct {
	logger *log.Logger
	coll   *mongo.Collection
}

func NewMemberGroupRepository(db *mongodb.MongoDB) MemberGroupRepository {
	coll := db.Client.
		Database(db.DBName).
		Collection(model.MemberGroupCollectionName)

	return &memberGroupRepo{
		logger: log.With("repository", "member_group_repo"),
		coll:   coll,
	}
}

// CountMembersByGroupID implements GroupRepository.
func (m *memberGroupRepo) CountMembersByGroupID(
	ctx context.Context,
	groupID string,
) (int64, error) {
	filter := bson.D{{Key: "group_id", Value: groupID}}
	return m.coll.CountDocuments(ctx, filter)
}

// FindMemberIDsByGroupID implements GroupRepository.
func (m *memberGroupRepo) FindMemberIDsByGroupID(
	ctx context.Context,
	groupID string,
) ([]string, error) {
	filter := bson.M{"group_id": groupID}

	cursor, err := m.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	memberGroups := []model.MemberGroup{}
	if err := cursor.All(ctx, &memberGroups); err != nil {
		return nil, err
	}

	ids := make([]string, len(memberGroups))
	for i, ele := range memberGroups {
		ids[i] = ele.MemberID
	}
	return ids, nil
}

// AddMembers implements MemberGroupRepository.
func (m *memberGroupRepo) AddMembers(
	ctx context.Context,
	groupID string,
	memberIDs []string,
) error {
	filter := bson.M{
		"group_id": groupID,
		"member_id": bson.M{
			"$in": memberIDs,
		},
	}

	cursor, err := m.coll.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	existingSet := make(map[string]struct{})
	for cursor.Next(ctx) {
		member := model.MemberGroup{}
		if err := cursor.Decode(&member); err != nil {
			return err
		}
		existingSet[member.MemberID] = struct{}{}
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	now := time.Now()
	newMembers := []model.MemberGroup{}
	for _, id := range memberIDs {
		if _, found := existingSet[id]; !found {
			newMembers = append(newMembers, model.MemberGroup{
				MemberID: id,
				GroupID:  groupID,
				BaseModel: model.BaseModel{
					ID:        uuid.NewString(),
					CreatedAt: now,
					UpdatedAt: now,
				},
			})
		}
	}

	if len(newMembers) == 0 {
		return nil
	}

	_, err = m.coll.InsertMany(ctx, newMembers)
	return err
}

// RemoveMembers implements MemberGroupRepository.
func (m *memberGroupRepo) RemoveMembers(
	ctx context.Context,
	groupID string,
	memberIDs []string,
) error {
	filter := bson.M{
		"group_id": groupID,
		"member_id": bson.M{
			"$in": memberIDs,
		},
	}

	_, err := m.coll.DeleteMany(ctx, filter)
	return err
}

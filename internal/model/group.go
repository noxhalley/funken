package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupStatus int8

const (
	GroupStatusActive GroupStatus = 1
	GroupStatusLocked GroupStatus = 2
)

type Group struct {
	BaseModel    `bson:",inline"       json:",inline"`
	Meta         bson.M               `bson:"meta,omitempty"         json:"meta"`
	Status       GroupStatus          `bson:"status"                 json:"status"`
	MemberIDs    []primitive.ObjectID `bson:"member_ids"             json:"member_ids"`
	MemberCount  *int                 `bson:"member_count,omitempty" json:"member_count,omitempty"`
	MessageCount int                  `bson:"message_count"          json:"message_count"`
}

func NewGroup(memberIDs []primitive.ObjectID) *Group {
	return &Group{
		BaseModel:    *NewBaseModel(),
		Meta:         bson.M{},
		Status:       GroupStatusActive,
		MemberCount:  nil,
		MessageCount: 0,
	}
}

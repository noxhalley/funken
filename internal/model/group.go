package model

import (
	"go.mongodb.org/mongo-driver/bson"
)

type GroupStatus int8

const (
	GroupStatusActive GroupStatus = 1
	GroupStatusLocked GroupStatus = 2

	GroupCollectionName = "groups"
)

type Group struct {
	BaseModel    `bson:",inline"       json:",inline"`
	Meta         bson.M      `bson:"meta,omitempty"         json:"meta"`
	Status       GroupStatus `bson:"status"                 json:"status"`
	MemberCount  *int        `bson:"member_count,omitempty" json:"member_count,omitempty"`
	MessageCount int         `bson:"message_count"          json:"message_count"`
}

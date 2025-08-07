package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MsgSortDirection string

const (
	MsgSortAsc  MsgSortDirection = "asc"
	MsgSortDesc MsgSortDirection = "desc"
)

type Message struct {
	BaseModel `bson:",inline"              json:",inline"`
	Message   string               `bson:"message"              json:"message"`
	GroupID   primitive.ObjectID   `bson:"group_id,omitempty"   json:"group_id"`
	SenderID  primitive.ObjectID   `bson:"sender_id,omitempty"  json:"sender_id"`
	Mentions  []primitive.ObjectID `bson:"mentions,omitempty"   json:"mentions"`
	Priority  bool                 `bson:"priority"             json:"priority"`
	Nickname  string               `bson:"nickname"             json:"nickname"`
	IPAddress string               `bson:"ip_address"           json:"ip_address"`
	DeletedAt *time.Time           `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

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
	BaseModel `bson:",inline"       json:",inline"`
	Message   string               `bson:"message"                       json:"message"`
	GroupID   primitive.ObjectID   `bson:"group_id,omitempty"            json:"groupId"`
	SenderID  primitive.ObjectID   `bson:"sender_id,omitempty"           json:"senderId"`
	Mentions  []primitive.ObjectID `bson:"mentions,omitempty"            json:"mentions"`
	Priority  bool                 `bson:"priority"                      json:"priority"`
	Nickname  string               `bson:"nickname"                      json:"nickname"`
	IPAddress string               `bson:"ip_address"                    json:"ipAddress"`
	DeletedAt *time.Time           `bson:"deleted_at,omitempty"          json:"deletedAt,omitempty"`
}

func NewMessage(
	message string,
	groupId primitive.ObjectID,
	senderId primitive.ObjectID,
	mentions []primitive.ObjectID,
	nickname string,
	ipAddress string,
	deletedAt *time.Time,
) *Message {
	return &Message{
		BaseModel: *NewBaseModel(),
		Message:   message,
		GroupID:   groupId,
		SenderID:  senderId,
		Mentions:  mentions,
		Priority:  false,
		Nickname:  nickname,
		IPAddress: ipAddress,
		DeletedAt: deletedAt,
	}
}

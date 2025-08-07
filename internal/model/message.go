package model

import "time"

type MsgSortDirection string

const (
	MsgSortAsc  MsgSortDirection = "asc"
	MsgSortDesc MsgSortDirection = "desc"
)

type Message struct {
	BaseModel `bson:",inline"              json:",inline"`
	Message   string     `bson:"message"              json:"message"`
	GroupID   string     `bson:"group_id,omitempty"   json:"group_id"`
	SenderID  string     `bson:"sender_id,omitempty"  json:"sender_id"`
	Mentions  []string   `bson:"mentions,omitempty"   json:"mentions"`
	Priority  bool       `bson:"priority"             json:"priority"`
	Nickname  string     `bson:"nickname"             json:"nickname"`
	IPAddress string     `bson:"ip_address"           json:"ip_address"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

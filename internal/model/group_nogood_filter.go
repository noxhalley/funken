package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type GroupNGFilter struct {
	BaseModel `bson:",inline"            json:",inline"`
	GroupID   primitive.ObjectID `bson:"group_id,omitempty" json:"groupId"`
	Title     string             `bson:"title"              json:"title"`
	Pattern   string             `bson:"pattern"            json:"pattern"`
	Flags     string             `bson:"flags,omitempty"    json:"flags,omitempty"`
}

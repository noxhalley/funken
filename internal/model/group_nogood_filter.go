package model

type GroupNGFilter struct {
	BaseModel `bson:",inline"            json:",inline"`
	GroupID   string `bson:"group_id,omitempty" json:"groupId"`
	Title     string `bson:"title"              json:"title"`
	Pattern   string `bson:"pattern"            json:"pattern"`
	Flags     string `bson:"flags,omitempty"    json:"flags,omitempty"`
}

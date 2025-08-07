package model

const MemberGroupCollectionName = "membergroups"

type MemberGroup struct {
	BaseModel `bson:",inline"       json:",inline"`
	MemberID  string `bson:"member_id"     json:"member_id"`
	GroupID   string `bson:"group_id"      json:"group_id"`
}

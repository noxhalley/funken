package model

import "time"

type BaseModel struct {
	ID        string    `bson:"id,omitempty"        json:"id"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at"`
}

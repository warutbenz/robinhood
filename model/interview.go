package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Interview struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Subject     string             `json:"subject" bson:"subject"`
	Detail      string             `json:"detail" bson:"detail"`
	Comments    []Comment          `json:"comments" bson:"comments"`
	Status      string             `json:"status" bson:"status"`
	CreateBy    string             `json:"create_by" bson:"create_by"`
	CreateDate  time.Time          `json:"create_date" bson:"create_date"`
	UpdatedBy   string             `json:"updated_by" bson:"updated_by"`
	UpdatedDate time.Time          `json:"updated_date" bson:"updated_date"`
}

type Comment struct {
	ID          int       `json:"id" bson:"id"`
	Comment     string    `json:"comment" bson:"comment"`
	CreateBy    string    `json:"create_by" bson:"create_by"`
	CreateDate  time.Time `json:"create_date" bson:"create_date"`
	UpdatedDate time.Time `json:"updated_date" bson:"updated_date"`
}

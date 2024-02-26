package database

import "go.mongodb.org/mongo-driver/bson/primitive"

// mongo models
type NodeInfo struct {
	OK        int                 `bson:"ok"`
	LocalTime *primitive.DateTime `bson:"localTime"`
}

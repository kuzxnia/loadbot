package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type AgentStatus struct {
	Name string `bson:"name"`

	// todo: add host

	// todo: add version

  // config?? 

	CreatedAt primitive.DateTime `bson:"created_at"`

	Heartbeat primitive.DateTime `bson:"heartbeat"`
}

package database

import (
	"github.com/kuzxnia/loadbot/lbot/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentStatus struct {
	Name string `bson:"name"`

	// todo: add host

	// todo: add version

	// config??

	CreatedAt primitive.DateTime `bson:"created_at"`

	Heartbeat primitive.DateTime `bson:"heartbeat"`
}

type Command struct {
	Id        primitive.ObjectID `bson:"_id"`
	Data      config.Job         `bson:"data"`
	Type      string             `bson:"state"`
	State     string             `bson:"state"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

type SubCommand struct {
	Id        primitive.ObjectID `bson:"_id"`
	Data      config.Job         `bson:"data"`
	Type      string             `bson:"state"`
	State     string             `bson:"state"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

// todo: move to different place
//
//go:generate stringer -type=CommandState -trimprefix=CommandState
type CommandState int

const (
	CommandStateCreated CommandState = iota // created
	CommandStateRunning                     // running
	CommandStateDone                        // done
	CommandStateError                       // error
)

//go:generate stringer -type=CommandType -trimprefix=CommandType
type CommandType int

const (
	CommandTypeStartWorkload CommandType = iota // start
	CommandTypeStopWorkload                     // stop
)

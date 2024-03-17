package database

import (
	"github.com/kuzxnia/loadbot/lbot/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentStatus struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Host      string             `bson:"host"`
  // add state
	CreatedAt primitive.DateTime `bson:"created_at"`
	Heartbeat primitive.DateTime `bson:"heartbeat"`
}

type Command struct {
	Id        primitive.ObjectID `bson:"_id"`
	Data      config.Job         `bson:"data"`
	Type      string             `bson:"type"`
	State     string             `bson:"state"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	Version   primitive.ObjectID `bson:"version"`
}

type Workload struct {
	Id        primitive.ObjectID `bson:"_id"`
	CommandId primitive.ObjectID `bson:"command_id"`
	AgentId   primitive.ObjectID `bson:"agent_id"`
	Data      config.Job         `bson:"data"`
	State     string             `bson:"state"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	Version   primitive.ObjectID `bson:"version"`
}

// todo: move to different place
//
//go:generate stringer -type=WorkloadState -trimprefix=WorkloadState
type WorkloadState int

const (
	WorkloadStateCreated WorkloadState = iota
	WorkloadStateToRun
	WorkloadStateRunning
	WorkloadStateDone
	WorkloadStateError
	WorkloadStateToDelete
	WorkloadStateDeleted
)

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

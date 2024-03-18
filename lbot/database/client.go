package database

import (
	"context"
	"time"

	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client interface {
	InsertOne(interface{}) (bool, error)
	InsertMany([]interface{}) (bool, error)
	ReadOne(interface{}) (bool, error)
	ReadMany(interface{}) (bool, error)
	UpdateOne(interface{}, interface{}) (bool, error)
	DropCollection() error
	Disconnect() error
}

type MongoClient struct {
	ctx        context.Context
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoClient(connectionString string, cfg *config.Job, schema *config.Schema) (*MongoClient, error) {
	opts := &options.ClientOptions{
		HTTPClient: HTTPClient(cfg),
	}
	opts = opts.
		ApplyURI(connectionString).
		SetReadPreference(readpref.SecondaryPreferred()).
		SetMaxPoolSize(cfg.Connections * 2).
		// SetMaxConnecting(100).
		SetMaxConnIdleTime(90 * time.Second).
		SetTimeout(cfg.Timeout)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		// log.Error("error in ping to mongo")
	} else {
		// log.Info("Successfully connected to database server")
	}

	var collection *mongo.Collection
	if schema != nil {
		collection = client.Database(schema.Database).Collection(schema.Collection)
	} else {
		collection = client.Database(cfg.Database).Collection(cfg.Collection)
	}
	return &MongoClient{ctx: ctx, client: client, collection: collection}, err
}

func NewInternalMongoClient(connectionString string) (*MongoClient, error) {
	opts := &options.ClientOptions{
		HTTPClient: HTTPClient(nil),
	}
	opts = opts.
		ApplyURI(connectionString).
		SetReadPreference(readpref.SecondaryPreferred()).
		SetMaxConnecting(100).
		SetMaxConnIdleTime(90 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())

	return &MongoClient{ctx: ctx, client: client, collection: nil}, err
}

func (c *MongoClient) Disconnect() (err error) {
	err = c.client.Disconnect(c.ctx)
	if err != nil {
		// log.Error("Error tring to disconnect from database", err)
	} else {
		// log.Info("Successfully disconnected from database server")
	}
	return
}

func (c *MongoClient) InsertOne(data interface{}) (bool, error) {
	_, err := c.collection.InsertOne(context.TODO(), data)
	return bool(err == nil), err
}

func (c *MongoClient) InsertMany(data []interface{}) (bool, error) {
	_, err := c.collection.InsertMany(context.TODO(), data)
	return bool(err == nil), err
}

func (c *MongoClient) ReadOne(filter interface{}) (bool, error) {
	var result bson.M
	err := c.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true, err
		}
		return false, err
	}
	return true, nil
}

func (c *MongoClient) ReadMany(filter interface{}) (bool, error) {
	batch_size := int32(1000)

	cursor, err := c.collection.Find(context.TODO(), bson.M{"author": "Franz Kafkaaa"}, &options.FindOptions{BatchSize: &batch_size})
	if err != nil {
		// log.Error(err)
	}

	defer cursor.Close(context.TODO())

	totalFound := 0
	for cursor.Next(context.Background()) {
		var data bson.M

		if err = cursor.Decode(&data); err != nil {
			// log.Error(err)
		}
		totalFound++
	}

	return true, nil
}

func (c *MongoClient) UpdateOne(filter interface{}, data interface{}) (bool, error) {
	// todo: only for now
	_, err := c.collection.UpdateOne(context.TODO(), filter, data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true, err
		}
		return false, err
	}
	return true, nil
}

func (c *MongoClient) DropCollection() error {
	return c.collection.Drop(context.TODO())
}

func (c *MongoClient) ClusterTime() (*primitive.DateTime, error) {
	res := c.client.Database(config.DB).RunCommand(context.TODO(), bson.D{{"isMaster", 1}})

	if err := res.Err(); err != nil {
		return nil, errors.WithMessage(err, "cmd: isMaster")
	}

	result := NodeInfo{}
	err := res.Decode(&result)

	return result.LocalTime, errors.WithMessage(err, "decode")
}

// workload
func (c *MongoClient) RunJob(job config.Job) error {
	// add lock??
	ct, err := c.ClusterTime()
	if err != nil {
		return errors.Wrap(err, "get cluster time")
	}

	// drop circular reference
	// job.Parent = nil

	cmd := Command{
		Id:        primitive.NewObjectID(),
		Data:      job,
		Type:      CommandTypeStartWorkload.String(),
		State:     CommandStateCreated.String(),
		CreatedAt: *ct,
	}

	_, err = c.client.Database(config.DB).Collection(config.CommandCollection).
		InsertOne(context.TODO(), cmd)

	if err != nil {
		return err
	}

	return nil
}

func (c *MongoClient) SaveWorkload(workload *Workload) error {
	oldVersion := workload.Version
	workload.Version = primitive.NewObjectID()

	_, err := c.client.Database(config.DB).Collection(config.WorkloadCollection).
		UpdateOne(context.TODO(), bson.M{"_id": workload.Id, "version": oldVersion}, bson.M{"$set": workload})
	if err != nil {
		workload.Version = oldVersion
		return err
	}

	return nil
}

func (c *MongoClient) AddWorkloads(workloads []interface{}) error {
	_, err := c.client.Database(config.DB).Collection(config.WorkloadCollection).
		InsertMany(context.TODO(), workloads)
	if err != nil {
		return err
	}

	return nil
}

func (c *MongoClient) SaveCommand(command *Command) error {
	oldVersion := command.Version
	command.Version = primitive.NewObjectID()

	_, err := c.client.Database(config.DB).Collection(config.CommandCollection).
		UpdateOne(context.TODO(), bson.M{"_id": command.Id, "version": oldVersion}, bson.M{"$set": command})
	if err != nil {
		// is this really necessary??
		command.Version = oldVersion
		return err
	}

	return nil
}

func (c *MongoClient) GetNextUnFinishedCommand() (*Command, error) {
	var cmd Command
	err := c.client.Database(config.DB).Collection(config.CommandCollection).
		FindOne(
			context.TODO(),
			bson.M{"state": bson.M{"$nin": bson.A{CommandStateDone.String(), CommandStateError.String()}}},
			&options.FindOneOptions{Sort: bson.M{"created_at": -1}},
		).Decode(&cmd)
	if err != nil {
		return nil, err
	}

	return &cmd, nil
}

// todo: temp change to stream for commands
func (c *MongoClient) GetNewWorkloads() (*Workload, error) {
	// add lock??
	var cmd Workload
	err := c.client.Database(config.DB).Collection(config.WorkloadCollection).
		FindOne(context.TODO(), bson.M{"state": WorkloadStateCreated.String()}).
		Decode(&cmd)
	if err != nil {
		return nil, err
	}

	return &cmd, nil
}

func (c *MongoClient) GetCommandWorkloads(command *Command) ([]*Workload, error) {
	cursor, err := c.client.Database(config.DB).Collection(config.WorkloadCollection).
		Find(context.TODO(), bson.M{"command_id": command.Id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	workloads := make([]*Workload, 0)

	for cursor.Next(context.Background()) {
		var workload Workload

		if err = cursor.Decode(&workload); err != nil {
			// log.Error(err)
		}
		workloads = append(workloads, &workload)

	}

	return workloads, nil
}

func (c *MongoClient) CancelCommand() error {
	return nil
}

// agent
func (c *MongoClient) GetAgentWithHeartbeatWithin() (uint64, error) {
	ct, err := c.ClusterTime()
	if err != nil {
		return 0, errors.Wrap(err, "get cluster time")
	}
	lastHbTime := ct.Time().Add(config.AgentsHeartbeatExpiration)

	cursor, err := c.client.Database(config.DB).Collection(config.AgentStatusCollection).Find(
		context.TODO(), bson.M{"heartbeat": bson.M{"$gte": lastHbTime}},
	)
	defer cursor.Close(context.TODO())

	totalFound := 0
	for cursor.Next(context.Background()) {
		var data bson.M

		if err = cursor.Decode(&data); err != nil {
			return 0, errors.Wrap(err, "error decoding agent")
		}
		totalFound++
	}
	return uint64(totalFound), nil
}

func (c *MongoClient) IsMasterAgent(agentId primitive.ObjectID) (bool, error) {
	ct, err := c.ClusterTime()
	if err != nil {
		return false, errors.Wrap(err, "get cluster time")
	}
	lastHbTime := primitive.NewDateTimeFromTime(ct.Time().Add(config.AgentsHeartbeatExpiration))

	var agent AgentStatus
	err = c.client.Database(config.DB).Collection(config.AgentStatusCollection).
		FindOne(context.TODO(), bson.M{"heartbeat": bson.M{"$gte": lastHbTime}}, &options.FindOneOptions{Sort: bson.M{"created_at": 1}}).
		Decode(&agent)

	if err != nil {
		return false, err
	}

	return agent.Id == agentId, nil
}

func (c *MongoClient) AddAgentStatus(stat AgentStatus) error {
	ct, err := c.ClusterTime()
	if err != nil {
		return errors.Wrap(err, "get cluster time")
	}
	stat.Heartbeat = *ct

	_, err = c.client.Database(config.DB).Collection(config.AgentStatusCollection).
		InsertOne(context.TODO(), stat)
	return errors.Wrap(err, "write into db")
}

func (c *MongoClient) SaveAgentStatus(stat AgentStatus) error {
	ct, err := c.ClusterTime()
	if err != nil {
		return errors.Wrap(err, "get cluster time")
	}
	stat.Heartbeat = *ct

	_, err = c.client.Database(config.DB).Collection(config.AgentStatusCollection).
		UpdateOne(context.TODO(), bson.M{"name": stat.Name, "_id": stat.Id}, bson.M{"$set": stat})
	return errors.Wrap(err, "write into db")
}

func (c *MongoClient) RemoveAgentStatus(stat AgentStatus) error {
	_, err := c.client.Database(config.DB).Collection(config.AgentStatusCollection).DeleteOne(
		context.TODO(), bson.M{"name": stat.Name},
	)
	return errors.WithMessage(err, "query")
}

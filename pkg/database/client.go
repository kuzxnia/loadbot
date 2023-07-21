package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kuzxnia/mongoload/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
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
		SetMaxConnecting(100).
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
		fmt.Println("error in ping to mongo")
	} else {
		fmt.Println("Successfully connected to database server")
	}

	collection := client.Database(schema.Database).Collection(schema.Collection)
	return &MongoClient{ctx: ctx, client: client, collection: collection}, err
}

func (c *MongoClient) Disconnect() error {
	return c.client.Disconnect(c.ctx)
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
	// start := time.Now()
	batch_size := int32(1000)

	cursor, err := c.collection.Find(context.TODO(), bson.M{"author": "Franz Kafkaaa"}, &options.FindOptions{BatchSize: &batch_size})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	// var results []Book
	// if err = cursor.All(context.Background(), &results); err != nil {
	// 	panic(err)
	// }

	//  println(len(results))

	// for _, book := range results {
	// 	res, _ := json.Marshal(book)
	// 	fmt.Println(string(res))
	// }

	totalFound := 0
	for cursor.Next(context.Background()) {
		var data bson.M

		if err = cursor.Decode(&data); err != nil {
			log.Fatal(err)
		}
		totalFound++
	}

	// elapsed := time.Since(start)
	// fmt.Printf("Find documents took %s", elapsed)
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

package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbClient struct {
	ctx        context.Context
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDbClient(uri string, databaseName string, collectionName string, maxPoolSize uint64) (*MongoDbClient, error) {
	if uri == "" {
		panic("uri is required")
	}

	opts := options.Client().
		ApplyURI(uri).
		SetReadPreference(readpref.SecondaryPreferred()).
		SetAppName("test").
		SetMaxPoolSize(maxPoolSize). // connectionsAmount * 8 is a magic number
		SetMaxConnecting(100).
		SetMaxConnIdleTime(time.Microsecond * 100000)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("connected")
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		fmt.Println("error in ping to mongo")
	} else {
		fmt.Println("no errors found")
	}

	// db - test, collection - go
	collection := client.Database(databaseName).Collection(collectionName)

	return &MongoDbClient{ctx: ctx, client: client, collection: collection}, err
}

func (c *MongoDbClient) Disconnect() error {
	return c.client.Disconnect(c.ctx)
}

func (c *MongoDbClient) InsertOne() (bool, error) {
	_, err := c.collection.InsertOne(context.Background(), SingleBook)

	return bool(err == nil), err
}

func (c *MongoDbClient) InsertMany() (bool, error) {
	_, err := c.collection.InsertMany(context.Background(), MultipleBooks)

	return bool(err == nil), err
}

func (c *MongoDbClient) readOne() (bool, error) {
	return true, nil
}

func (c *MongoDbClient) readMany() (bool, error) {
	// start := time.Now()
	batch_size := int32(1000)

	cursor, err := c.collection.Find(context.Background(), bson.M{"author": "Franz Kafkaaa"}, &options.FindOptions{BatchSize: &batch_size})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

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
		var book Book

		if err = cursor.Decode(&book); err != nil {
			log.Fatal(err)
		}
		totalFound++
	}

	// elapsed := time.Since(start)
	// fmt.Printf("Find documents took %s", elapsed)
	return true, nil
}

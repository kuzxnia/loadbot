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
	ctx           context.Context
	client        *mongo.Client
	collection    *mongo.Collection
	batchProvider *DataProvider
}

// todo: change params to options struct
func NewMongoDbClient(
	uri string,
	databaseName string,
	collectionName string,
	maxPoolSize uint64,
	batchSize uint64,
	dataLenght uint64,
) (*MongoDbClient, error) {
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
	}
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		fmt.Println("error in ping to mongo")
	} else {
		fmt.Println("Successfully connected to database server")
	}

	// db - test, collection - go
	collection := client.Database(databaseName).Collection(collectionName)

	batchProvider := NewDataProvider(batchSize, dataLenght)

	return &MongoDbClient{ctx: ctx, client: client, collection: collection, batchProvider: batchProvider}, err
}

func (c *MongoDbClient) Disconnect() error {
	return c.client.Disconnect(c.ctx)
}

func (c *MongoDbClient) InsertOne() (bool, error) {
	_, err := c.collection.InsertOne(context.Background(), c.batchProvider.singleItem)
	return bool(err == nil), err
}

func (c *MongoDbClient) InsertMany() (bool, error) {
	_, err := c.collection.InsertMany(context.Background(), *c.batchProvider.batchOfItems)
	return bool(err == nil), err
}

func (c *MongoDbClient) InsertOneOrMany() (bool, error) {
	err := error(nil)
	if c.batchProvider.batchSize == 0 {
		_, err = c.collection.InsertOne(context.Background(), c.batchProvider.singleItem)
	} else {
		_, err = c.collection.InsertMany(context.Background(), *c.batchProvider.batchOfItems)
	}

	return bool(err == nil), err
}

func (c *MongoDbClient) ReadOne() (bool, error) {
	return true, nil
}

func (c *MongoDbClient) ReadMany() (bool, error) {
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

func (c *MongoDbClient) GetBatchSize() uint64 {
	return c.batchProvider.batchSize
}

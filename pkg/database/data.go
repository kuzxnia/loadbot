package database

import (
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type DataProvider struct {
	batchSize uint64

	singleItem   *bson.M
	batchOfItems *[]interface{}
}

func NewDataProvider(batchSize, dataLenght uint64) *DataProvider {
	if dataLenght < 0 {
		panic("dataLenght must be positive")
	}
	if batchSize < 0 {
		panic("batchSize must be positive")
	}

	batchOfData := make([]interface{}, batchSize)

	for i := 0; i < int(batchSize); i++ {
		batchOfData[i] = bson.M{"data": randStringBytes(dataLenght)}
	}

	return &DataProvider{
		batchSize:    batchSize,
		batchOfItems: &batchOfData,
		singleItem:   &bson.M{"data": randStringBytes(dataLenght)},
	}
}

func randStringBytes(n uint64) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

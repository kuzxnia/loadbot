package schema

import (
	"math/rand"

	"github.com/kuzxnia/mongoload/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
)

type DataGenerator interface {
	Generate(*config.Schema, uint64) interface{}
}

func NewDataGenerator() DataGenerator {
	return DataGenerator(&FakeDataGenerator{})
}

type FakeDataGenerator struct{}

func (g *FakeDataGenerator) Generate(schema *config.Schema, dataSize uint64) interface{} {
	// todo: use faker to generate map'like data
	// check size of empty bson to calculate how much data generate
	return &bson.M{"data": randStringBytes(dataSize)}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n uint64) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

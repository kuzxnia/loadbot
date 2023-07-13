package schema

import (
	"github.com/kuzxnia/mongoload/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
)

type DataProvider interface {
	GetSingleItem() interface{}
	GetBatch(uint64) []interface{}
	GetFilter() interface{}
}

func NewDataProvider(schema *config.Schema) DataProvider {
	return DataProvider(NewSimpleDataProvider(100, 100))
}

// here i need to create pool of items to be taken to insert/update
// also here is place to store keys which needs to be saved

type SimpleDataProvider struct {
	singleItem   *bson.M
	batchOfItems *[]interface{}
}

func NewSimpleDataProvider(batchSize, dataLenght uint64) *SimpleDataProvider {
	generator := NewDataGenerator()
	batchOfData := make([]interface{}, batchSize)

	for i := 0; i < int(batchSize); i++ {
		batchOfData[i] = generator.Generate(nil, dataLenght)
	}

	return &SimpleDataProvider{
		batchOfItems: &batchOfData,
		singleItem:   &bson.M{"data": randStringBytes(dataLenght)},
	}
}

func (d *SimpleDataProvider) GetSingleItem() interface{} {
	return d.singleItem
}

func (d *SimpleDataProvider) GetFilter() interface{} {
  // todo: 
	return d.singleItem
}

func (d *SimpleDataProvider) GetBatch(batchSize uint64) []interface{} {
	// todo: add slice
	return *d.batchOfItems
}

package schema

import (
	"github.com/kuzxnia/mongoload/pkg/config"
)

type DataProvider interface {
	GetSingleItem() interface{}
	GetBatch(uint64) []interface{}
	GetFilter() interface{}
}

func NewDataProvider(job *config.Job) DataProvider {
	return DataProvider(
		NewLiveDataProvider(job.BatchSize, job.DataSize, job.GetSchema()),
	)
}

// here i need to create pool of items to be taken to insert/update
// also here is place to store keys which needs to be saved

type LiveDataProvider struct {
	dataGenerator DataGenerator
}

// todo: generate on file and take from pool
// type PoolDataProvider struct { }
func NewLiveDataProvider(batchSize, dataSize uint64, schema *config.Schema) *LiveDataProvider {
	return &LiveDataProvider{
		dataGenerator: NewDataGenerator(schema, dataSize),
	}
}

func (d *LiveDataProvider) GetSingleItem() interface{} {
	singleItem, _ := d.dataGenerator.Generate()
	return singleItem
}

func (d *LiveDataProvider) GetFilter() interface{} {
	// todo:
	return d.GetSingleItem()
}

func (d *LiveDataProvider) GetBatch(batchSize uint64) []interface{} {
	batchOfData := make([]interface{}, batchSize)

	for i := 0; i < int(batchSize); i++ {
		batchOfData[i], _ = d.dataGenerator.Generate()
	}

	// todo: add slice
	return batchOfData
}

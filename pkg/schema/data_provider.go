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
		NewSimpleDataProvider(job.BatchSize, job.DataSize, job.GetTemplateSchema()),
	)
}

// here i need to create pool of items to be taken to insert/update
// also here is place to store keys which needs to be saved

type SimpleDataProvider struct {
	dataGenerator DataGenerator
	singleItem    *interface{}
	batchOfItems  *[]interface{}
}

// generate upfront - better for workload test not for leading data
func NewSimpleDataProvider(batchSize, dataSize uint64, schema *config.Schema) *SimpleDataProvider {
	generator := NewDataGenerator(schema, dataSize)
	batchOfData := make([]interface{}, batchSize)

	for i := 0; i < int(batchSize); i++ {
		batchOfData[i], _ = generator.Generate()
	}

	return &SimpleDataProvider{
		dataGenerator: generator,
		batchOfItems:  &batchOfData,
		singleItem:    &batchOfData[0],
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

// todo: generate on file and take from pool
// type PoolDataProvider struct { }

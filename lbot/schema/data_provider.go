package schema

import (
	"github.com/kuzxnia/loadbot/lbot/config"
)

type DataProvider interface {
	GetSingleItem() interface{}
	GetSingleItemWithout(string) interface{}
	GetBatch(uint64) []interface{}
	GetFilter() interface{}
}

func NewDataProvider(job *config.Job, schema *config.Schema) DataProvider {
	return DataProvider(
		NewLiveDataProvider(job, schema),
	)
}

// here i need to create pool of items to be taken to insert/update
// also here is place to store keys which needs to be saved

type LiveDataProvider struct {
	job           *config.Job
	dataGenerator DataGenerator
}

// todo: generate on file and take from pool
// type PoolDataProvider struct { }
func NewLiveDataProvider(job *config.Job, schema *config.Schema) *LiveDataProvider {
	return &LiveDataProvider{
		job:           job,
		dataGenerator: NewDataGenerator(schema, job.DataSize),
	}
}

func (d *LiveDataProvider) GetSingleItem() interface{} {
	singleItem, _ := d.dataGenerator.Generate()
	return singleItem
}

// todo: remote this, add abstraction with skipped keys or sth like that
func (d *LiveDataProvider) GetSingleItemWithout(key string) interface{} {
	singleItem, _ := d.dataGenerator.Generate()
	v := singleItem.(map[string]interface{})
	delete(v, key)
	return v
}

func (d *LiveDataProvider) GetFilter() interface{} {
	singleItem, _ := d.dataGenerator.GenerateFromTemplate(d.job.Filter)
	return singleItem
}

func (d *LiveDataProvider) GetBatch(batchSize uint64) []interface{} {
	batchOfData := make([]interface{}, batchSize)

	for i := 0; i < int(batchSize); i++ {
		batchOfData[i], _ = d.dataGenerator.Generate()
	}

	// todo: add slice
	return batchOfData
}

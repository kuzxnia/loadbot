package driver

import (
	"time"

	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/database"
	"github.com/kuzxnia/loadbot/lbot/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type JobHandler interface {
	Execute() error
}

func NewJobHandler(job *config.Job, client database.Client, dataPool schema.DataPool) JobHandler {
	dataProvider := schema.NewDataProvider(job)
	handler := BaseHandler{
		job:          job,
		client:       client,
		dataProvider: dataProvider,
		dataPool:     dataPool,
	}

	switch job.Type {
	case string(config.Write):
		return JobHandler(&WriteHandler{BaseHandler: &handler})
	case string(config.Read):
		return JobHandler(&ReadHandler{BaseHandler: &handler})
	case string(config.Update):
		return JobHandler(&UpdateHandler{BaseHandler: &handler})
	case string(config.BulkWrite):
		return JobHandler(&BulkWriteHandler{BaseHandler: &handler})
	case string(config.DropCollection):
		return JobHandler(&DropCollection{BaseHandler: &handler})
	case string(config.Sleep):
		return JobHandler(&SleepHandler{Duration: job.Duration})
	default:
		// todo change
		panic("Invalid job type: " + job.Type)
	}
}

type BaseHandler struct {
	job          *config.Job
	client       database.Client
	dataProvider schema.DataProvider
	dataPool     schema.DataPool
}

type WriteHandler struct {
	*BaseHandler
}

func (h *WriteHandler) Execute() error {
	item := h.dataProvider.GetSingleItem()

	_, error := h.client.InsertOne(item)

	if error == nil && h.dataPool != nil {
		h.dataPool.Set(item)
	}
	return error
}

type BulkWriteHandler struct {
	*BaseHandler
}

func (h *BulkWriteHandler) Execute() error {
	items := h.dataProvider.GetBatch(100)

	_, error := h.client.InsertMany(items)

	if error == nil && h.dataPool != nil {
		h.dataPool.SetBatch(items)
	}
	return error
}

type ReadHandler struct {
	*BaseHandler
}

func (h *ReadHandler) Execute() error {
	filter := h.dataProvider.GetFilter()

	_, error := h.client.ReadOne(filter)
	return error
}

type UpdateHandler struct {
	*BaseHandler
}

func (h *UpdateHandler) Execute() error {
	item := h.dataProvider.GetSingleItemWithout("_id")
	filter := h.dataProvider.GetFilter()

	_, error := h.client.UpdateOne(filter, bson.M{"$set": item})
	return error
}

type DropCollection struct {
	*BaseHandler
}

func (h *DropCollection) Execute() error {
	error := h.client.DropCollection()
	return error
}

type SleepHandler struct {
	Duration time.Duration
}

func (h *SleepHandler) Execute() error {
	time.Sleep(h.Duration)
	return nil
}

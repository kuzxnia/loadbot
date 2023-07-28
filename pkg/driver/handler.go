package driver

import (
	"time"

	"github.com/kuzxnia/mongoload/pkg/config"
	"github.com/kuzxnia/mongoload/pkg/database"
	"github.com/kuzxnia/mongoload/pkg/schema"
)

type JobHandler interface {
	Handle() (time.Duration, error)
}

func NewJobHandler(job *config.Job, client database.Client, dataPool schema.DataPool) JobHandler {
	handler := BaseHandler{
		job:          job,
		client:       client,
		dataProvider: schema.NewDataProvider(job),
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
	case string(config.CreateIndex):
		return JobHandler(&CreateIndex{BaseHandler: &handler})
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

func (h *WriteHandler) Handle() (time.Duration, error) {
	item := h.dataProvider.GetSingleItem()

	start := time.Now()
	_, error := h.client.InsertOne(item)
	elapsed := time.Since(start)

	if error == nil && h.dataPool != nil {
		h.dataPool.Set(item)
	}
	return elapsed, error
}

type BulkWriteHandler struct {
	*BaseHandler
}

func (h *BulkWriteHandler) Handle() (time.Duration, error) {
	items := h.dataProvider.GetBatch(100)

	start := time.Now()
	_, error := h.client.InsertMany(items)
	elapsed := time.Since(start)

	if error == nil && h.dataPool != nil {
		h.dataPool.SetBatch(items)
	}
	return elapsed, error
}

type ReadHandler struct {
	*BaseHandler
}

func (h *ReadHandler) Handle() (time.Duration, error) {
	start := time.Now()
	_, error := h.client.ReadOne(h.dataProvider.GetFilter())
	elapsed := time.Since(start)
	return elapsed, error
}

type UpdateHandler struct {
	*BaseHandler
}

func (h *UpdateHandler) Handle() (time.Duration, error) {
	start := time.Now()
	_, error := h.client.UpdateOne(h.dataProvider.GetFilter(), h.dataProvider.GetSingleItem())
	elapsed := time.Since(start)
	return elapsed, error
}

type CreateIndex struct {
	*BaseHandler
}

func (h *CreateIndex) Handle() (time.Duration, error) {
	start := time.Now()
	error := h.client.CreateIndexes(h.job.Indexes)
	elapsed := time.Since(start)
	return elapsed, error
}

type DropCollection struct {
	*BaseHandler
}

func (h *DropCollection) Handle() (time.Duration, error) {
	start := time.Now()
	error := h.client.DropCollection()
	elapsed := time.Since(start)
	return elapsed, error
}

type SleepHandler struct {
	Duration time.Duration
}

func (h *SleepHandler) Handle() (time.Duration, error) {
	time.Sleep(h.Duration)
	return h.Duration, nil
}

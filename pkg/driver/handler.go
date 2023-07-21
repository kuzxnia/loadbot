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

func NewJobHandler(job *config.Job, client database.Client) JobHandler {
	// todo: move provider to outside of this to use generated data in all workers
	handler := BaseHandler{
		client:   client,
		provider: schema.NewDataProvider(job),
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
	client   database.Client
	provider schema.DataProvider
}

type WriteHandler struct {
	*BaseHandler
}

func (h *WriteHandler) Handle() (time.Duration, error) {
	start := time.Now()
	_, error := h.client.InsertOne(h.provider.GetSingleItem())
	elapsed := time.Since(start)
	return elapsed, error
}

type BulkWriteHandler struct {
	*BaseHandler
}

func (h *BulkWriteHandler) Handle() (time.Duration, error) {
	start := time.Now()
	_, error := h.client.InsertMany(h.provider.GetBatch(100))
	elapsed := time.Since(start)
	return elapsed, error
}

type ReadHandler struct {
	*BaseHandler
}

func (h *ReadHandler) Handle() (time.Duration, error) {
	start := time.Now()
	_, error := h.client.ReadOne(h.provider.GetFilter())
	elapsed := time.Since(start)
	return elapsed, error
}

type UpdateHandler struct {
	*BaseHandler
}

func (h *UpdateHandler) Handle() (time.Duration, error) {
	start := time.Now()
	_, error := h.client.UpdateOne(h.provider.GetFilter(), h.provider.GetSingleItem())
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

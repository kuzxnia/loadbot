package schema

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/kuzxnia/loadbot/lbot/config"
)

type DataPool interface {
	Get(string) interface{}
	Set(interface{})
	SetBatch([]interface{})
	ExtendGeneratorMapperFields(generator *GeneratorFieldMapper)
}

func NewDataPool(schema *config.Schema) DataPool {
	return DataPool(
		&InMemoryDataPool{
			schema: schema,
			dataPool: make(map[string]struct {
				pointer uint64
				data    []interface{}
			}),
		},
	)
}

type InMemoryDataPool struct {
	schema   *config.Schema
	mutex    sync.RWMutex
	dataPool map[string]struct {
		pointer uint64
		data    []interface{}
	}
}

func (d *InMemoryDataPool) Get(key string) (result interface{}) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if dataPool, ok := d.dataPool[key]; ok {
		pointer := atomic.LoadUint64(&dataPool.pointer)
		if dataLen := len(dataPool.data); dataLen == 0 {
			return
		} else if dataLen <= int(pointer) {
			atomic.StoreUint64(&dataPool.pointer, 0)
		}
		result = dataPool.data[pointer]
		atomic.AddUint64(&dataPool.pointer, 1)
	}
	return
}

func (d *InMemoryDataPool) Set(data interface{}) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, key := range d.schema.Save {
		dataPool, ok := d.dataPool[key]
		if !ok {
			d.dataPool[key] = struct {
				pointer uint64
				data    []interface{}
			}{}
		}
		field, err := GetFieldFromData(key, data)
		if err == nil {
			dataPool.data = append(dataPool.data, field)
		}
		d.dataPool[key] = dataPool
	}
}

func (d *InMemoryDataPool) SetBatch(dataBatch []interface{}) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, key := range d.schema.Save {
		dataPool, ok := d.dataPool[key]
		if !ok {
			d.dataPool[key] = struct {
				pointer uint64
				data    []interface{}
			}{}
		}
		for _, data := range dataBatch {
			field, err := GetFieldFromData(key, data)
			if err == nil {
				dataPool.data = append(dataPool.data, field)
			}
		}
		d.dataPool[key] = dataPool
	}
}

func (d *InMemoryDataPool) ExtendGeneratorMapperFields(generator *GeneratorFieldMapper) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for key := range d.dataPool {
		// todo: fix problem with same field returned, propably related to go closures
		generator.Set("#"+key, func(opts ...options.OptionFunc) string {
			return d.Get(key).(string)
		})
	}
}

func GetFieldFromData(fieldPath string, rawData interface{}) (data interface{}, err error) {
	data = rawData
	for _, key := range strings.Split(fieldPath, ".") {
		v, ok := data.(map[string]interface{})
		if !ok {
			return nil, errors.New("Data not contains saved field")
		}
		data, ok = v[key]
		if !ok {
			return nil, errors.New("Data not contains saved field")
		}
	}
	return
}

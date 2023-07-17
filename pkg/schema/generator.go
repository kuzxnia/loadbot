package schema

import (
	"errors"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/kuzxnia/mongoload/pkg/config"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

var GeneratorFieldTypesMapper = map[string]func(opts ...options.OptionFunc) string{
	"#string": faker.Word,
	"#word":   faker.Word,
	"#id":     faker.UUIDDigit,
}

var GeneratorFieldTypes = lo.Keys(GeneratorFieldTypesMapper)

type DataGenerator interface {
	Generate() (interface{}, error)
}

func NewDataGenerator(schema *config.Schema, dataSize uint64) DataGenerator {
	if schema != nil {
		return DataGenerator(
			&StructuralizableDataGenerator{
				schema:   schema,
				dataSize: dataSize,
			},
		)
	}
	return DataGenerator(
		&SimpleDataGenerator{
			dataSize: dataSize,
		},
	)
}

type SimpleDataGenerator struct {
	dataSize uint64
}

func (g *SimpleDataGenerator) Generate() (interface{}, error) {
	// todo: use faker to generate map'like data
	// check size of empty bson to calculate how much data generate
	return &bson.M{"data": randStringBytes(g.dataSize)}, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n uint64) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type StructuralizableDataGenerator struct {
	schema   *config.Schema
	dataSize uint64
}

func (g *StructuralizableDataGenerator) Generate() (interface{}, error) {
	return g.generate(g.schema.Schema)
}

// recurent func for parsing with building nested bson
func (g *StructuralizableDataGenerator) generate(value interface{}) (interface{}, error) {
	switch value := value.(type) {
	case string:
		if generate, ok := GeneratorFieldTypesMapper[value]; ok {
			return generate(), nil
		} else {
			return nil, errors.New("Invalid field mapper, got: " + value)
		}

	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, v := range value {
			value, err := g.generate(v)
			if err != nil {
				return nil, err
			}
			result[k] = value
		}
		return result, nil
	default:
		return nil, errors.New("Invalid schema format")
	}
}

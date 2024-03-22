package schema

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	DefaultGeneratorFieldMapper = NewGeneratorFieldMapper()
	// todo: better validation
	GeneratorFieldTypes = lo.Keys(DefaultGeneratorFieldMapper.FieldTypeMapper)
)

// todo: add interface
type GeneratorFieldMapper struct {
	// todo: make private
	FieldTypeMapper map[string]func(opts ...options.OptionFunc) string
}

func NewGeneratorFieldMapper() *GeneratorFieldMapper {
	return &GeneratorFieldMapper{
		FieldTypeMapper: map[string]func(opts ...options.OptionFunc) string{
			"#id":     faker.UUIDDigit,
			"#string": faker.Word,
			"#word":   faker.Word,
			// internet
			"#email":    faker.Email,
			"#username": faker.Username,
			"#password": faker.Password,
			// person
			"#name":              faker.Name,
			"#first_name":        faker.FirstName,
			"#first_name_male":   faker.FirstNameMale,
			"#first_name_female": faker.FirstNameFemale,
			"#last_name":         faker.LastName,
			"#title_male":        faker.TitleMale,
			"#title_female":      faker.TitleFemale,
			"#phone_number":      faker.Phonenumber,
		},
	}
}

func (m *GeneratorFieldMapper) Generate(field string) (result interface{}, err error) {
	if generate, ok := m.FieldTypeMapper[field]; ok {
		return generate(), nil
	} else {
		return nil, errors.New("Invalid field mapper, got: " + field)
	}
}

func (m *GeneratorFieldMapper) Set(field string, valueGenerator func(opts ...options.OptionFunc) string) {
	m.FieldTypeMapper[field] = valueGenerator
}

type DataGenerator interface {
	Generate() (interface{}, error)
	GenerateFromTemplate(interface{}) (interface{}, error)
}

func NewDataGenerator(schema *config.Schema, dataSize int) DataGenerator {
	// todo: check size of object using, unsafe.Sizeof( )

	if schema != nil {
		return DataGenerator(
			&StructuralizableDataGenerator{
				schema: schema,
				// add support for custom byte size
			},
		)
	}
	return DataGenerator(
		&MeasurableDataGenerator{
			dataSize: dataSize,
		},
	)
}

type MeasurableDataGenerator struct {
	dataSize int
}

func (g *MeasurableDataGenerator) Generate() (interface{}, error) {
	return &bson.M{"data": randStringBytes(g.dataSize)}, nil
}

func (g *MeasurableDataGenerator) GenerateFromTemplate(template interface{}) (interface{}, error) {
	switch value := template.(type) {
	case string:
		return randStringBytes(g.dataSize), nil

	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, nestedTemplate := range value {
			value, err := g.GenerateFromTemplate(nestedTemplate)
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

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randStringBytes(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

type StructuralizableDataGenerator struct {
	schema *config.Schema
}

func (g *StructuralizableDataGenerator) Generate() (interface{}, error) {
	result, error := g.GenerateFromTemplate(g.schema.Schema)
	return result, error
}

// recurent func for parsing with building nested bson
func (g *StructuralizableDataGenerator) GenerateFromTemplate(template interface{}) (interface{}, error) {
	switch value := template.(type) {
	case string:
		generatedValue, err := DefaultGeneratorFieldMapper.Generate(value)
		if err != nil {
			return nil, errors.New("Invalid field mapper, got: " + value)
		}
		return generatedValue, nil

	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, nestedTemplate := range value {
			value, err := g.GenerateFromTemplate(nestedTemplate)
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

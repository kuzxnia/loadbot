package schema

import (
	"testing"

	"github.com/kuzxnia/mongoload/pkg/config"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDataFromSchema(t *testing.T) {
	schema := config.Schema{
		Name: "dummy",
		Schema: map[string]interface{}{
			"name":     "#string",
			"surname":  "#string",
			"lastname": "#string",
			"address": map[string]interface{}{
				"street":   "#string",
				"postcode": "#string",
			},
		},
	}
	generator := NewDataGenerator(&schema, 100)

	result, error := generator.Generate()

	assert.Nil(t, error)
	assert.Subset(t, lo.Keys(schema.Schema), lo.Keys(result.(map[string]interface{})))
	assert.NotEmpty(t, lo.Values(result.(map[string]interface{})))
}

func TestInvalidType(t *testing.T) {
	schema := config.Schema{
		Name: "dummy",
		Schema: map[string]interface{}{
			"name": "#invalid",
		},
	}
	generator := NewDataGenerator(&schema, 100)

	result, error := generator.Generate()

	assert.Nil(t, result)
	assert.Error(t, error, "Invalid field mapper, got: #invalid")
}

package config


import "errors"

func (c *Config) Validate() error {
	validators := []func() error{
		c.validateAllJobTemplatesAreProvided,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (c *Config) validateAllJobTemplatesAreProvided() error {
	isSchemaName := func(schema *Schema, comparator string) bool {
		return schema.Name == comparator
	}

	for _, job := range c.Jobs {
		if !Contains[*Schema, string](c.Schemas, job.Template, isSchemaName) {
			return errors.New("Job: " + job.Name + " have invalid template \"" + job.Template + "\"")
		}
	}
	return nil
}

// todo: validation job type
// todo: validation duration and opertions cannot be set together

func Contains[T comparable, X comparable](array []T, comparator X, predicate func(T, X) bool) bool {
	for _, elem := range array {
		if predicate(elem, comparator) {
			return true
		}
	}
	return false
}

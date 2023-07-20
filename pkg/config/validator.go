package config

import "errors"

func (c *Config) Validate() error {
	validators := []func() error{
		c.validateJobs,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (c *Config) validateJobs() error {
	for _, job := range c.Jobs {
		if error := job.Validate(); error != nil {
			return error
		}
	}
	return nil
}

func (job *Job) Validate() error {
	validators := []func() error{
		job.validateTemplate,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (job *Job) validateTemplate() error {
	isSchemaName := func(schema *Schema, comparator string) bool {
		return schema.Name == comparator
	}

	if !Contains[*Schema, string](job.Parent.Schemas, job.Template, isSchemaName) {
		return errors.New("Job: " + job.Name + " have invalid template \"" + job.Template + "\"")
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

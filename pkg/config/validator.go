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
		job.validateType,
		job.validateDuration,
		job.validatePace,
		job.validateConnections,
		job.validateBatchSize,
		job.validateOperations,
		job.validateDataSize,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (job *Job) validateTemplate() error {
	if string(Sleep) == job.Type {
		return nil
	}

	isSchemaName := func(schema *Schema, comparator string) bool {
		return schema.Name == comparator
	}

	if !Contains[*Schema, string](job.Parent.Schemas, job.Template, isSchemaName) {
		return errors.New("Job: " + job.Name + " have invalid template \"" + job.Template + "\"")
	}
	return nil
}

func (job *Job) validateType() (err error) {
	switch job.Type {
	case string(Write):
	case string(BulkWrite):
	case string(Read):
	case string(Update):
	case string(DropCollection):
	case string(Sleep):
	default:
		err = errors.New("Job type: " + job.Type + " ")
	}
	return
}

func (job *Job) validateConnections() (err error) {
	if job.Connections == 0 {
		err = errors.New("JobValidationError: field 'connections' must be greater than 0")
	}
	if job.Type == string(Sleep) {
		if job.Connections != 1 {
			err = errors.New("JobValidationError: field 'connections' max number concurrent connections for job type 'sleep' is 1")
		}
	}
	return
}

func (job *Job) validateDuration() (err error) {
	if job.Type == string(Sleep) {
		if job.Duration <= 0 {
			err = errors.New("JobValidationError: field 'duration' must be greater than 0 for job with 'sleep' type ")
		}
	}
	return
}

func (job *Job) validatePace() (err error) {
	if job.Type == string(Sleep) {
		if job.Pace != 0 {
			err = errors.New("JobValidationError: field 'pace' must be equal 0 or must be not set for job with 'sleep' type ")
		}
	}
	return
}

func (job *Job) validateBatchSize() (err error) {
	if job.Type == string(Sleep) {
		if job.BatchSize != 0 {
			err = errors.New("JobValidationError: field 'batch_size' must be equal 0 or must be not set for job with 'sleep' type ")
		}
	}
	return
}

func (job *Job) validateDataSize() (err error) {
	if job.Type == string(Sleep) {
		if job.DataSize != 0 {
			err = errors.New("JobValidationError: field 'data_size' must be equal 0 or must be not set for job with 'sleep' type ")
		}
	}
	return
}

func (job *Job) validateOperations() (err error) {
	if job.Type == string(Sleep) {
		if job.Operations != 0 {
			err = errors.New("JobValidationError: field 'operations' must be equal 0 or must be not set for job with 'sleep' type ")
		}
	}
	return
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
package config

import "errors"

func (c *Config) Validate() error {
	validators := []func() error{
		c.validateJobs,
		// c.validateSchemas,
		c.validateReportingFormats,
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

func (c *Config) validateReportingFormats() error {
	for _, reportingFormat := range c.ReportingFormats {
		if error := reportingFormat.Validate(); error != nil {
			return error
		}
	}
	return nil
}

func (job *Job) Validate() error {
	validators := []func() error{
		job.validateSchema,
		job.validateReportFormat,
		job.validateDatabase,
		job.validateCollection,
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

func (job *Job) validateSchema() error {
	if string(Sleep) == job.Type || job.Schema == "" {
		return nil
	}

	if !Contains(job.Parent.Schemas, func(s *Schema) bool { return s.Name == job.Schema }) {
		return errors.New("JobValidationError: job \"" + job.Name + "\" have invalid template \"" + job.Schema + "\"")
	}
	return nil
}

func (job *Job) validateReportFormat() error {
	if job.ReportingFormat == "" {
		return nil
	}

	if !Contains(job.Parent.ReportingFormats, func(s *ReportingFormat) bool { return s.Name == job.ReportingFormat }) {
		return errors.New("JobValidationError: job \"" + job.Name + "\" have invalid report_format \"" + job.ReportingFormat + "\"")
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

func (job *Job) validateDatabase() (err error) {
	if job.Schema != "" || job.Type == string(Sleep) {
		return
	}
	if job.Database == "" {
		err = errors.New("JobValidationError: field 'database' is required if 'template' or 'type' is not set")
	}
	return
}

func (job *Job) validateCollection() (err error) {
	if job.Schema != "" || job.Type == string(Sleep) {
		return
	}
	if job.Collection == "" {
		err = errors.New("JobValidationError: field 'collection' is required if 'template' or 'type' is not set")
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

func (rp *ReportingFormat) Validate() error {
	validators := []func() error{
		rp.validateReportingFormat,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (rpt *ReportingFormat) validateReportingFormat() (err error) {
	return nil
}

// todo: validation job type
// todo: validation duration and opertions cannot be set together

// func Contains[T comparable, X comparable](array []T, comparator X, predicate func(T, X) bool) bool {
// 	for _, elem := range array {
// 		if predicate(elem, comparator) {
// 			return true
// 		}
// 	}
// 	return false
// }

func Contains[T comparable](array []T, predicate func(T) bool) bool {
	for _, elem := range array {
		if predicate(elem) {
			return true
		}
	}
	return false
}

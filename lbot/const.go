package lbot

import "time"

type JobType string

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"
)

const (
	Write          JobType = "write"
	BulkWrite      JobType = "bulk_write"
	Read           JobType = "read"
	Update         JobType = "update"
	Sleep          JobType = "sleep"
	DropCollection JobType = "drop_collection"
)

var DefaultReportFormats = []*ReportingFormat{
	{
		Name:     "default",
		Interval: 15 * time.Second,
		Template: `{{.Now}} {{if .JobName -}}Job: "{{.JobName}}" {{else -}}Job type: "{{.JobType}}"{{end}}
Reqs: {{.TotalReqs}}, RPS {{f2 .Rps}}, s:{{.SuccessReqs}}/err:{{.ErrorReqs}}/tout:{{.TimeoutErr}}/errRate:{{f1 .ErrorRate}}%
AVG: {{msf3 .Avg}}ms P50: {{msf3 .P50}}ms, P90: {{msf3 .P90}}ms P99: {{msf3 .P99}}ms

`,
	},
	{
		Name:     "simple",
		Interval: 15 * time.Second,
		Template: "{{.Now}} Reqs: {{.TotalReqs}}, RPS {{f2 .Rps}} s:{{.SuccessReqs}}/err:{{.ErrorReqs}}\n\n",
	},
	{
		Name:     "write",
		Interval: 15 * time.Second,
		Template: `{{.Now}} Reqs: {{.TotalReqs}}, RPS {{f2 .Rps}}, s:{{.SuccessReqs}}/err:{{.ErrorReqs}}/tout:{{.TimeoutErr}}/errRate:{{f1 .ErrorRate}}%
AVG: {{msf3 .Avg}}ms P50: {{msf3 .P50}}ms, P90: {{msf3 .P90}}ms P99: {{msf3 .P99}}ms

`,
	},
	{
		Name:     "bulk_write",
		Interval: 15 * time.Second,
		Template: `{{.Now}} Reqs: {{.TotalReqs}}, OPS: {{.TotalOps}}, RPS {{f2 .Rps}}, OPS {{f2 .Ops}}, s:{{.SuccessReqs}}/err{{.ErrorReqs}}/tout:{{.TimeoutErr}}/errRate:{{f1 .ErrorRate}}
AVG: {{msf3 .Avg}}ms P50: {{msf3 .P50}}ms, P90: {{msf3 .P90}}ms P99: {{msf3 .P99}}ms

`,
	},
}
var DefaultReportFormat = DefaultReportFormats[0]

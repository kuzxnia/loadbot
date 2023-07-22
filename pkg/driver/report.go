package driver

import (
	"fmt"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/kuzxnia/mongoload/pkg/config"
	"github.com/montanaflynn/stats"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ReportFormatKeys            = NewReportFormatKeys()
	DefaultReportFormatTemplate = `{{if .JobName -}} Job: "{{.JobName}}" {{else -}} Job type: "{{.JobType}}"{{end}}
Total reqs: {{.TotalReqs}}, RPS {{f2 .Rps}} success: {{.SuccessReqs}}, errors: {{.ErrorReqs}} timeout: {{.TimeoutErr}}, error rate: {{f1 .ErrorRate}}
AVG: {{f3 .Avg}}s P50: {{f3 .P50}}s, P90: {{f3 .P90}}s P99: {{f3 .P99}}s

`
)

func NewReportFormatKeys() []string {
	keys := []string{
		"{{.SuccessReqs}}",
		"{{.TotalReqs}}",
		"{{.ErrorReqs}}",
		"{{.TimeoutErr}}",
		"{{.NoDataErr}}",
		"{{.OtherErr}}",
		"{{.ErrorRate}}",
		"{{.Min}}",
		"{{.Max}}",
		"{{.Avg}}",
		"{{.Rps}}",
    // if batch := w.db.GetBatchSize(); batch > 1 {
    // 	fmt.Printf("Operations per second: %f op/s\n", float64(requestsDone*batch)/elapsed.Seconds())
    // }
		// ops
	}
	for i := 1; i < 100; i++ {
		keys = append(keys, fmt.Sprintf("{{.P%v}}", i))
	}
	return keys
}

type Report interface {
	Len() int
	Sum() (float64, error)
	Min() (float64, error)
	Max() (float64, error)
	Mean() (float64, error)
	Percentile(float64) (float64, error)
	Percentiles(input ...float64) (percentiles []float64, err error)

	Add(time.Duration, error)
	SetDuration(time.Duration)
	Summary()
}

func NewReport(job *config.Job) Report {
	return Report(
		&TemplateReport{
			job:  job,
			data: make([]time.Duration, 0),
		},
	)
}

type TemplateReport struct {
	job                   *config.Job
	mutex                 sync.RWMutex
	data                  []time.Duration
	rawData               []float64
	errorsReqs            uint64
	timeoutErrors         uint64
	noDocumentsFoundError uint64
	duration              time.Duration
}

func (s *TemplateReport) Len() int               { return len(s.data) }
func (s *TemplateReport) Sum() (float64, error)  { return stats.Sum(*s.GetRawData()) }
func (s *TemplateReport) Min() (float64, error)  { return stats.Min(*s.GetRawData()) }
func (s *TemplateReport) Max() (float64, error)  { return stats.Max(*s.GetRawData()) }
func (s *TemplateReport) Mean() (float64, error) { return stats.Mean(*s.GetRawData()) }
func (s *TemplateReport) Percentile(percentile float64) (float64, error) {
	return stats.Percentile(*s.GetRawData(), percentile)
}

func (s *TemplateReport) Percentiles(input ...float64) (percentiles []float64, err error) {
	percentiles = make([]float64, len(input))
	for i, percentile := range input {
		percentiles[i], err = stats.Percentile(*s.GetRawData(), percentile)
	}
	return
}

func (s *TemplateReport) GetRawData() *[]float64 {
	// need to lock this for
	// readlock s.data, lock rawData
	if len(s.data) != len(s.rawData) {
		for i := len(s.rawData); i < len(s.data); i++ {
			s.rawData = append(s.rawData, float64(s.data[i].Seconds()))
		}
	}
	return &s.rawData
}

func (s *TemplateReport) SetDuration(duration time.Duration) {
	s.duration = duration
}

func (s *TemplateReport) GetReportData() map[string]any {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// handle errors
	totalReqs := s.Len()
	min, _ := s.Min()
	max, _ := s.Max()
	avg, _ := s.Mean()

	mapping := map[string]any{
		"JobName":     s.job.Name,
		"JobType":     s.job.Type,
		"SuccessReqs": totalReqs - int(s.errorsReqs),
		"TotalReqs":   totalReqs,
		"ErrorReqs":   s.errorsReqs,
		"TimeoutErr":  s.timeoutErrors,
		"noDataErr":   s.noDocumentsFoundError,
		"OtherErr":    s.errorsReqs - s.timeoutErrors - s.noDocumentsFoundError,
		"ErrorRate":   IfElse(totalReqs != 0, float64(s.errorsReqs)/float64(totalReqs)*100, 0),
		"Min":         min,
		"Max":         max,
		"Avg":         avg,
		"Rps":         IfElse(s.duration != 0, float64(totalReqs)/float64(s.duration.Seconds()), 0),
	}
	var key string
	for i := 1; i < 100; i++ {
		key = fmt.Sprintf("P%v", i)
		p, _ := s.Percentile(float64(i))
		mapping[key] = p
	}
	return mapping
}

func (s *TemplateReport) Add(interval time.Duration, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = append(s.data, interval)
	if err != nil {
		s.errorsReqs++
		if mongo.IsTimeout(err) {
			s.timeoutErrors++
		} else if err == mongo.ErrNoDocuments {
			s.noDocumentsFoundError++
		}
	}
}

func (s *TemplateReport) Summary() {
	// todo: handle error, by default Must panics
	outputTemplate := template.Must(template.New("").Funcs(template.FuncMap{
		"f1": func(f float64) string { return fmt.Sprintf("%.1f", f) },
		"f2": func(f float64) string { return fmt.Sprintf("%.2f", f) },
		"f3": func(f float64) string { return fmt.Sprintf("%.3f", f) },
		"f4": func(f float64) string { return fmt.Sprintf("%.4f", f) },
	}).Parse(DefaultReportFormatTemplate))
	outputTemplate.Execute(os.Stdout, s.GetReportData())
}

func IfElse[T comparable](condition bool, a T, b T) T {
	if condition {
		return a
	}
	return b
}

// var ErrNoDocuments = errors.New("mongo: no documents in result")
// type ActorHistogram struct {
// 	Report
// 	datach chan float64
// }

// func NewActorHistogram() Histogram {
// 	data := make([]float64, 0)
// 	hist := &ActorHistogram{
// 		BaseHistogram: BaseHistogram{
// 			data: data,
// 		},
//     datach: make(chan float64),
// 	}

// 	go hist.loop()

// 	return Histogram(hist)
// }

// func (h *ActorHistogram) Update(interval float64) {
// 	h.datach <- interval
// }

// func (h *ActorHistogram) loop() {
// 	for interval := range h.datach {
// 		h.data = append(h.data, interval)
// 	}
// }

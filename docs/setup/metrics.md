The metrics exporting feature enables you to monitor and analyze various metrics related to your workload. Below, we detail the key aspects of this feature:

### Library Dependency
Our project utilizes the lightweight library [github.com/VictoriaMetrics/metrics](https://github.com/VictoriaMetrics/metrics) as an alternative to [github.com/prometheus/client_golang](https://github.com/prometheus/client_golang). This library facilitates efficient handling and exporting of metrics for monitoring purposes.

### Exporting Strategies
You have the flexibility to choose between two strategies for exporting metrics: pull or push.

#### Push Strategy:

With this strategy, the tool pushes metrics to the provided URL at the specified interval.

Configure the [agent](/loadbot/setup/agent/) using the following flags:

- **metrics_export_url**: Specifies the URL where metrics will be pushed.
- **metrics_export_interval_seconds**: Defines the interval (in seconds) for pushing metrics.

#### Pull Strategy:
With this strategy, external app pull metrics from us and put them into Prometheus.

Configure the [agent](/loadbot/setup/agent/) using the following flags:
- **metrics_export_url**: Configure the agent to expose an HTTP server on the provided port.


### Metrics Selection

The metrics available for monitoring encompass both system-level metrics and custom metrics related to your workloads. Here are some popular metrics from a Go application's perspective:

##### System Metrics:

- `go_goroutines`
- `go_threads`
- `go_memstats_alloc_bytes`
- `go_memstats_heap_inuse_bytes`
- `go_memstats_alloc_bytes_total`
- `process_resident_memory_bytes`

##### Custom Workload Metrics:

- `requests_total`
- `requests_error`
- `requests_duration_seconds`

#### Labels for Querying
When querying custom workload metrics, you can utilize labels to specify job-related information:

```
{job="job_name_here", job_uuid="auto_generate_uuid", job_type="write|bulk_write|read|update..." agent="agent_name_here"}
```

##### Example query:
```
requests_total{job="workload 1", agent="186.12.9.19"}
```

The `job_uuid` label distinguishes between different job runs/attempts, allowing you to track and analyze performance across multiple executions of the same job. Additionally, all metrics are labeled with the `name of the agent`, enabling you to differentiate metrics coming from different agents.

### Additional Resources
Metrics have been extracted from VictoriaMetrics sources. For more in-depth information about VictoriaMetrics, you can refer to the following article: [VictoriaMetrics: Creating the Best Remote Storage for Prometheus](https://faun.pub/victoriametrics-creating-the-best-remote-storage-for-prometheus-5d92d66787ac).

You can also refer to an excellent article discussing Prometheus metrics in Go, which can provide further insights (Prometheus Go Metrics)[https://povilasv.me/prometheus-go-metrics/].

Example Grafana Golang dashboard [template](https://grafana.com/grafana/dashboards/10376-generic-go-process/).

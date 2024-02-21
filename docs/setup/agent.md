# Agent configuration

The agent configuration allows you to customize various parameters for the behavior of the agent. Below is an example file illustrating the structure and available fields for configuring the agent:
> Note:
> You can configure the agent both from the command line interface and from a configuration file.

### CLI Usage
To start the loadbot agent, use the following command:

    Usage:
      loadbot start-agent [flags]

    Flags:
      -f, --config-file string                     Config file for loadbot-agent
      -h, --help                                   help for start-agent
          --metrics_export_interval_seconds uint   Prometheus export push interval
          --metrics_export_port string             Expose metrics on port instead pushing to prometheus
          --metrics_export_url string              Prometheus export url used for pushing metrics
      -n, --name string                            Agent name
      -p, --port string                            Agent port
          --stdin                                  Provide configuration from stdin.

> Note:
> Configurations specified from the command line interface will overwrite those from the configuration file.

### Configuration file

In the configuration file, you're able to configure both workloads and the agent, whereas via the command-line interface, you can solely configure the agent.

```json
{
    "agent": {
        "name": "mongo 6.0.8 workload",
        "port": "1234",
        "metrics_export_url": "http://victoria-metrics:8428/api/v1/import/prometheus",
        "metrics_export_interval_seconds": 10,
        "metrics_export_port": "9090",
    }
}
```

> Note: 
> The agent configuration is only applied when the agent is started. It cannot be changed at runtime using commands such as `loadbot config`.

### Agent Fields
- **name** (string, optional): Specifies the name of the agent. This field is used as a label in the metrics exporter. Useful when exporting metrics from multiple agents.
- **port** (string, optional): Specifies the port on which the agent listens for incoming connections.
- **metrics_export_url** (string, optional): If set, metrics will be pushed to the specified endpoint in the Prometheus standard format.
- **metrics_export_interval_seconds** (integer, optional): Specifies the interval (in seconds) at which metrics are pushed to the export endpoint.
- **metrics_export_port** (string, optional): If set, metrics will be accessible on the provided port.



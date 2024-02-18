
After installing loadbot, you can quickly get started by following these steps:

1. Run loadbot agent with your desired configuration using the command:
```bash
loadbot start-agent -f config_file.json
```

2. Start the workload using the loadbot client:
```bash
loadbot start
```

3. Monitor the progress of the workload using the command:
```bash
loadbot progress
```

4. To stop the workload, use the following command:
```bash
loadbot stop
```

### Example config file

For more information on defining [jobs](https://kuzxnia.github.io/loadbot/loadbot/setup/job/) and [schemas](https://kuzxnia.github.io/loadbot/loadbot/setup/schema/), visit the [configuration](https://kuzxnia.github.io/loadbot/loadbot/setup/) section of the documentation.

```json
{
  "connection_string": "mongodb://myadmin:abc123@127.0.0.1:27017",
  "jobs": [
    {  // start on empty collection - fresh start 
      "type": "drop_collection",
      "database": "tmp-kuzxnia-showcase-db",
      "collection": "tmp-kuzxnia-showcase-col",
      "operations": 1
    },
    {
      "name": "Workload 20s max throughput",
      "type": "write",
      "database": "tmp-kuzxnia-showcase-db",
      "collection": "tmp-kuzxnia-showcase-col",
      "data_size": 200,
      "connections": 100,
      "duration": "20s",
    },
    {  // give little rest to db
      "type": "sleep",
      "duration": "5s",
      "format": "simple"
    },
    {
      "name": "Workload 1_000_000 elements with pace 10K rps",
      "type": "write",
      "batch_size": 100,
      "database": "tmp-kuzxnia-showcase-db",
      "collection": "tmp-kuzxnia-showcase-col",
      "data_size": 100,
      "connections": 200,
      "pace": 10000,
      "operations": 1000000,
    }
}

```

### CLI usage

```
$ loadbot
A command-line database workload driver

Usage:
  lbot [command]

Agent Commands:
  start-agent Start lbot-agent

Driver Commands:
  config      Config
  progress    Watch stress test
  start       Start stress test
  stop        Stopping stress test

Additional Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -u, --agent-uri string    loadbot agent uri (default: 127.0.0.1:1234) (default "127.0.0.1:1234")
  -h, --help                help for lbot
      --log-format string   log format, must be one of: json, fancy (default "fancy")
      --log-level string    log level, must be one of: trace, debug, info, warn, error, fatal, panic (default "info")
  -v, --version             version for lbot

Use "lbot [command] --help" for more information about a command.
```

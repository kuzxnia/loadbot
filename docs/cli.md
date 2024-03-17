---
hide:
  - navigation
---





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

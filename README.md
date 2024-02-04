# lbot - workload database driver (currenty supporting only MongoDB)

## Introduction
The purpose of this tool is to simulate workloads to facilitate testing the failover capabilities of database cluster under load. This code, being an open-source project, is in its early development stage and likely contains various bugs.


## How to use:
1. Build image - `make build-docker`
2. Run agent - `make run-docker-agent`
3. Configure agent by running lbot config command - `make run-docker-lbot-config CONFIG_FILE="_config.json"`
4. Run workload tests - `make run-docker-lbot start`

> Note: If running with local db remember to use host network and configure connection_string to 127.0.0.1 `docker run --network="host" --rm -t mload < config_file.json` or check Makefile

This tool offers two ways to access it: one through CLI arguments and the other via a configuration file. Utilizing the configuration file provides additional functionalities for the tool.

### CLI usage:
    A command-line database workload

    Usage:
      lbot [command]

    Driver Commands:
      config      Config
      start       Start stress test
      stop        Stop stress test

    Additional Commands:
      completion  Generate the autocompletion script for the specified shell
      help        Help about any command

    Flags:
      -u, --agent-uri string    loadbot agent uri (default: 127.0.0.1:1234)
      -h, --help                help for lbot
      -v, --version             version for lbot

    Use "lbot [command] --help" for more information about a command.


Known issue:
* srv not working with some DNS servers - golang 1.13+ issue see [this](https://github.com/golang/go/issues/37362) and [this](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#hdr-Potential_DNS_Issues)

    > Old versions of kube-dns and the native DNS resolver (systemd-resolver) on Ubuntu 18.04 are known to be non-compliant in this manner. 

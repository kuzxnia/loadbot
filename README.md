# mongoload - workload driver for MongoDB 

## Introduction
The purpose of this tool is to simulate workloads to facilitate testing the failover capabilities of MongoDB clusters under load. This code, being an open-source project, is in its early development stage and likely contains various bugs.


## How to use:
1. Build image - `docker build -t mload .`
2. Run - `docker run mload -c 10 -o 10000 http://localhost:21017`

This tool offers two ways to access it: one through CLI arguments and the other via a configuration file. Utilizing the configuration file provides additional functionalities for the tool.

### CLI usage:
    Arguments:
      [<connection-string>]     Database connection string

    Flags:
    -h, --help                  Show context-sensitive help.
    -c, --connections=10        Number of concurrent connections
    -p, --pace=UINT-64          Pace - RPS limit
    -d, --duration=DURATION     Duration (ex. 10s, 5m, 1h)
    -o, --operations=UINT-64    Operations (read/write/update) to perform
    -b, --batch-size=UINT-64    Batch size
    -t, --timeout=5s            Timeout for requests
    -f, --config-file=STRING    Config file path
        --debug                 Displaying additional diagnostic information


### Configuration file:
Due to the limited functionalities of the CLI, in order to fully harness the capabilities of this tool, it is advisable to utilize a configuration file. The program can be executed by specifying the configuration file with `--config-file <file-path>` or `-f <file-path>`. For instance, the command `docker run mload -f path/to/config/file.json` demonstrates how to use this approach.

Example file:

```json
{
  "connection_string": "mongodb://localhost:27017",
  "debug": true,
  "jobs": [
    {
      "name": "default job",
      "type": "write",
      "template": "default",
      "connections": 100,
      "pace": 0,
      "data_size": 0,
      "batch_size": 0,
      "duration": "0s",
      "operations": 1000,
      "timeout": "1s"
    }
  ],
  "schemas": [
    {
      "name": "default",
      "database": "load_test",
      "collection": "load_test",
      "schema": {
        "_id": "#_id",
        "name": "#string",
        "lastname": "#string"
      }
    }
  ]
}
```
<details>
<summary>Defining schemas</summary>

</br>
**Schema fields**

- `name` - unique name, used in jobs (see job.schema) for determining which template use
- `database` - database name
- `collection` - collection name
- `schema` - actual document template

**Schema document template fields:**

General
- `#id` 
- `#string`
- `#word`

Internet
- `#email`
- `#username`
- `#password`
 
Person
- `#name`
- `#first_name`
- `#first_name_male`
- `#first_name_female`
- `#last_name`
- `#title_male`
- `#title_female`
- `#phone_number`

**More examples**

</details>

<details>
<summary>Defining Jobs</summary>
</details>

<br>
**Simple job example**

```json
{
  "jobs": [
    {
      "name": "default job",
      "type": "write",
      "template": "default",
      "connections": 100,
      "pace": 0,
      "data_size": 0,
      "batch_size": 0,
      "duration": "0s",
      "operations": 1000,
      "timeout": "1s"
    }
  ],
}
```

**Jobs fields:**

* `name`(string) - job name
* `type`(enum `write|bulk_write|read|update`) - operation type
* `template`(string) - schema name, if you will not provide schema data will be inserted in `{'data': <generate_data>}` format
* `connection`(unsigned int) - number of concurrent connections, number is not limited to physical threads number
* `data_size`(unsigned int) - data size inserted (currently only works for default schema)
* `batch_size`(unsigned int) - insert batch size (only applicable for `bulk_write` job type)
* `duration`(string) - duration time ex. 1h, 15m, 10s
* `operations`(unsigned int) - number of requests to perform, ex. 100 reads, 100 bulk_writes
* `timeout`(string) - connection timeout ex. 1h, 15m, 10s


> Note:
> If you don't provide the requests amount or duration limit program will continue running 
> indefinitely unless it is manually stopped by pressing `ctrl-c`. 


Known issue:
* srv not working with k8s 1.17 - it is golang 1.13+ issue see (this)[https://github.com/golang/go/issues/37362]



# mongoload - workload driver for MongoDB 

## Introduction
The purpose of this tool is to simulate workloads to facilitate testing the failover capabilities of MongoDB clusters under load. This code, being an open-source project, is in its early development stage and likely contains various bugs.


## How to use:
1. Build image - `docker build -t mload .`
2. Run with simple cli aruments - `docker run -it mload -c 10 -o 10000 mongodb://lovely.host.name:27017`
3. Run with config file - `docker run -it mload < config_file.json`

> Note: If running with local db remember to use host network and configure connection_string to 127.0.0.1 `docker run --network="host" --rm -it mload < config_file.json`

This tool offers two ways to access it: one through CLI arguments and the other via a configuration file. Utilizing the configuration file provides additional functionalities for the tool.

### CLI usage:
    Arguments:
      [<connection-string>]     Database connection string

    Flags:
    -h, --help                  Show context-sensitive help.
    -c, --connections=1         Number of concurrent connections
    -p, --pace=UINT-64          Pace - RPS limit
    -d, --duration=DURATION     Duration (ex. 10s, 5m, 1h)
    -o, --operations=UINT-64    Operations (read/write/update) to perform
    -b, --batch-size=UINT-64    Batch size
    -t, --timeout=5s            Timeout for requests
    -f, --config-file=STRING    Config file path
        --debug                 Displaying additional diagnostic information


### Configuration file:
Due to the limited functionalities of the CLI, in order to fully harness the capabilities of this tool, it is advisable to utilize a configuration file. The program can be executed by specifying the configuration file with:
1. stdin with binary `cat config_file.json | mload` or stdin with docker `docker run mload -it mload < config_file.json`
2. flag with `--config-file <file-path>` or `-f <file-path>`. 

> Note: If you want to use config file with docker you need to mount volume with file or copy when building image.

Example file:

```json
{
  "connection_string": "mongodb://localhost:27017",
  "debug": true,
  "jobs": [
    {
      "name": "Write 100c 1k ops",
      "type": "write",
      "schema": "user_schema",
      "connections": 100,
      "operations": 1000,
    },
    {
      "name": "Dummy job name/ read 30s 100rps",
      "type": "read",
      "schema": "user_schema",
      "connections": 100,
      "pace": 100,
      "duration": "30s",
      "filter": {
          "special_name": "#special_name"
      }
    }
  ],
  "schemas": [
    {
      "name": "user_schema",
      "database": "load_test",
      "collection": "load_test",
      "schema": {
        "_id": "#_id",
        "special_name": "#string",
        "lastname": "#string"
      },
      "save": [
        "special_name"  // will be avaliable in job.filter under "#special_name"
      ]
    },
  ],
  "reporting_formats": [
    {
      "name": "simple",
      "interval": "5s",
      "template": "Job: {{.JobType}}, total reqs: {{.TotalReqs}}, RPS {{f2 .Rps}} success: {{.SuccessReqs}}\n\n"
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

<br>

**Example write with schema 100ops**

```json
{
  "name": "insert with schema",
  "type": "write",
  "schema": "user_schema",
  "connections": 10,
  "operations": 100
}
```

**Write without schema 20s**

```json
{
  "name": "insert without schema",
  "type": "write",
  "database": "load_test",
  "collection": "load_test",
  "connections": 10,
  "data_size": 100,
  "duration": "20s",
  "timeout": "1s"
}
```

**Read with schema 20s**

```json
{
  "name": "read with schema",
  "type": "read",
  "schema": "user_schema",
  "connections": 10,
  "operations": 100,
  "filter": {
    "user_name": "#user_name",
    "name": "#generate_value"  // here you can use remember/saved value as well as generated one
  }
}
```

**Let the database rest**

```json
{
  "type": "sleep",
  "duration": "5s"
}
```

**Drop collection**

```json
{
  "type": "drop_collection",
  "database": "load_test",
  "collection": "load_test",
  "operations": 1
}
```
or with schema
```json
{
  "type": "drop_collection",
  "schema": "example_schema",
  "operations": 1
}
```

**Jobs fields:**

* `name`(string, optional) - job name
* `type`(enum `write|bulk_write|read|update|create_index|drop_collection|sleep`) - operation type
* `template`(string) - schema name, if you will not provide schema data will be inserted in `{'data': <generate_data>}` format
* `database`(string, required if schema is not set) - database name
* `schema`(string, optional) - string foreign-key to schemas list
* `filter`(string, required for read and update) - filter schema
* `indexes`(list, optional) - list of indexes to create (only for type "create_index") 
* `format`(string, optional) - string foreign-key to reporting_formats list
* `collection`(string, required if schema is not set) - collection name
* `connection`(unsigned int) - number of concurrent connections, number is not limited to physical threads number
* `data_size`(unsigned int) - data size inserted (currently only works for default schema)
* `batch_size`(unsigned int) - insert batch size (only applicable for `bulk_write` job type)
* `duration`(string) - duration time ex. 1h, 15m, 10s
* `operations`(unsigned int) - number of requests to perform, ex. 100 reads, 100 bulk_writes
* `timeout`(string) - connection timeout ex. 1h, 15m, 10s

</details>

<details>
<summary>Custom reporting format</summary>

<br>

By default, you have access to several available output formats: `default`, `simple`, `write`, and `bulk_write`. If you do not provide the 'default' format, the default format associated with the type will be utilized. Moreover, you have the flexibility to override all formats except for the `default` one.

**default**

```
2023/08/12 15:10:44 Job: "lovely job name"
Reqs: 7500, RPS 500.00, s:7500/err:0/tout:0/errRate:0.0%
AVG: 1.395ms P50: 0.484ms, P90: 0.703ms P99: 8.048ms
```

**simple**

```
2023/08/12 17:28:47 Reqs: 2592, RPS 500.05 s:2592/err:0
```

**write**

```
2023/08/12 17:28:52 Reqs: 2595, RPS 499.91, s:2595/err:0/tout:0/errRate:0.0%
AVG: 1.987ms P50: 0.989ms, P90: 1.634ms P99: 34.556ms

```

**bulk_write**

```
2023/08/12 17:28:57 Reqs: 2590, OPS: 259000, RPS 499.50, OPS 49949.52, s:2590/err0/tout:0/errRate:0.0
AVG: 6.867ms P50: 3.814ms, P90: 6.786ms P99: 79.938ms
```

**Example custom reporting format**

```json
{
  "name": "custom",
  "interval": "5s",
  "template": "{{.Now}} Job: {{.JobType}}, total reqs: {{.TotalReqs}}, RPS {{f2 .Rps}} success: {{.SuccessReqs}}\n\n"
}
```
- `name` - used to determine which template to use (see section job.format)
- `interval` - if set, tests reports/summaries will be displayed at set time intervals
- `template` - report format


**Template fields**

`Now`, `JobName`, `JobType`, `JobBatchSize`,`SuccessReqs`, `ErrorReqs`, `TotalReqs`, `TotalOps`, `TimeoutErr`, `NoDataErr`, `OtherErr`, `ErrorRate`, `Rps`, `Ops`

**Math fields**

`Min`, `Max`, `Avg`, `Rps` and `P<number>` ex. `P90` - percentiles

**Floating point fields formatters**

`f<number>` - format number to n places (1 to 4) ex. `{{f2 .Rps}}` 

`msf<number>` - format number to n places (1 to 4) and convert to milliseconds ex. `{{msf2 .P99}}` 


</details>

<details>
<summary>More examples</summary>

<br>

- Index creation job
```json
{
  "type": "create_index",
  "template": "default",
  "indexes": [
    {
      "keys": { "name": 1 },
      "options": { "unique": false, "name": "dummy_name_index_name" },
    }
  ]
}
```
or without using schema
```json
{
  "type": "create_index",
  "database": "load_test",
  "collection": "load_test",
  "operations": 1,
  "indexes": [
    {
      "keys": {"name": 1},
    }
  ]
}
```

</details>

<details>
<summary>Other features</summary>

<br>

**Features**

- JSON standardization - comments and trailing commas support ex.
```json
{
    "jobs": [
        {
          "type": "drop_collection",
          "database": "load_test",
          "collection": "load_test",
          "operations": 1
        },
        /*{
          "type": "sleep",
          "duration": "5s",
          "format": "simple"
        },*/
    ]
}
```


</details>

> Note:
> If you don't provide the requests amount or duration limit program will continue running 
> indefinitely unless it is manually stopped by pressing `ctrl-c`. 


Known issue:
* srv not working with some DNS servers - golang 1.13+ issue see [this](https://github.com/golang/go/issues/37362) and [this](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#hdr-Potential_DNS_Issues)

    > Old versions of kube-dns and the native DNS resolver (systemd-resolver) on Ubuntu 18.04 are known to be non-compliant in this manner. 



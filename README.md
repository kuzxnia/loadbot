# mongoload - workload driver for MongoDB 

## Introduction
The purpose of this tool is to simulate workloads to facilitate testing the failover capabilities of MongoDB clusters under load. This code, being an open-source project, is in its early development stage and likely contains various bugs.


## How to use:
1. Build image - `docker build -t mload .`
2. Run - `docker run mload -uri=http://localhost:21017 -req=10000`

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


### Simple config-file example:
You can execute the program using `--config-file <file-path>` or `-f <file-path>`. The file should be in JSON format. 
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


**Schema fields**

- `name` - unique name, used in jobs (see job.schema) for determining which template use
- `database` - database name
- `collection` - collection name
- `schema` - actual document template

**Schema document template fields:**

General
- `#id` - ex. ``
- `#string` - ex. ``
- `#word`` - ex. ``

Internet
- `#email`` - ex. ``
- `#username`` - ex. ``
- `#password`` - ex. ``
 
Person
- `#name`` - ex. ``
- `#first_name`` - ex. ``
- `#first_name_male`` - ex. ``
- `#first_name_female`` - ex. ``
- `#last_name`` - ex. ``
- `#title_male`` - ex. ``
- `#title_female`` - ex. ``
- `#phone_number`` - ex. ``

**More examples**

</details>

<details>
<summary>Defining Jobs</summary>
</details>



> Note:
> If you don't provide the requests amount or duration limit program will continue running 
> indefinitely unless it is manually stopped by pressing `ctrl-c`. 

Limitations:


Known issue:



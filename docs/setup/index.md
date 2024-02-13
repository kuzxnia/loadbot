# Configuration file

Due to the limited functionalities of the CLI, in order to fully harness the capabilities of this tool, it is advisable to utilize a configuration file.

> Note: If you want to use config file with docker you need to mount volume with config file.

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


After installing loadbot, you can quickly get started by following these steps:

1. Run loadbot agent with your desired configuration using the command:
```bash
loadbot start-agent -f config_file.json
```

2. Start and watch the workload using the loadbot client:
```bash
loadbot start --progress
Job "My first job" |██████████████████████████████████████████████████████████████████| 30/30S 50RPS 1509REQ
```

3. To stop the workload, use the following command:
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

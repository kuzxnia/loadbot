
### Jobs fields:

- `name`(string, optional) - job name
- `type`(enum `write|bulk_write|read|update|create_index|drop_collection|sleep`) - operation type
- `template`(string) - schema name, if you will not provide schema data will be inserted in `{'data': <generate_data>}` format
- `database`(string, required if schema is not set) - database name
- `schema`(string, optional) - string foreign-key to schemas list
- `filter`(string, required for read and update) - filter schema
- `indexes`(list, optional) - list of indexes to create (only for type "create_index") 
- `format`(string, optional) - string foreign-key to reporting_formats list
- `collection`(string, required if schema is not set) - collection name
- `connection`(unsigned int) - number of concurrent connections, number is not limited to physical threads number
- `data_size`(unsigned int) - data size inserted (currently only works for default schema)
- `batch_size`(unsigned int) - insert batch size (only applicable for `bulk_write` job type)
- `duration`(string) - duration time ex. 1h, 15m, 10s
- `operations`(unsigned int) - number of requests to perform, ex. 100 reads, 100 bulk_writes
- `timeout`(string) - connection timeout ex. 1h, 15m, 10s


### Defining Jobs
Example write with schema 100ops

```json
{
  "name": "insert with schema",
  "type": "write",
  "schema": "user_schema",
  "connections": 10,
  "operations": 100
}
```

### Write without schema 20s

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

### Read with schema 20s

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

### Let the database rest

```json
{
  "type": "sleep",
  "duration": "5s"
}
```

### Drop collection

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

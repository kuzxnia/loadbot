# Other features

- config file JSON standardization - comments and trailing commas support ex.
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

> Note:
> If you don't provide the requests amount or duration limit program will continue running 
> indefinitely unless it is manually stopped by pressing `ctrl-c`. 


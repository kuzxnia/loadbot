# mongoload - workload driver for MongoDB 

## Introduction
The purpose of this tool is to simulate workloads to facilitate testing the failover capabilities of MongoDB clusters under load. This code, being an open-source project, is in its early development stage and likely contains various bugs.


## How to use:
1. Build image - `docker build -t mload .`
2. Run - `docker run mload -uri=http://localhost:21017 -req=10000`

### Current flags:
    -h, --help                                 Show context-sensitive help.
    -u, --uri="mongodb://localhost:27017"      Database hostname url
        --mongo-database="load_test"           Database name
        --mongo-collection="load_test_coll"    Collection name
    -c, --concurrent-connections=100           Concurrent connections amount
        --rps=INT                              RPS limit
    -d, --duration=DURATION                    Duration limit
    -r, --requests=INT                         Requests to perform
    -b, --batch-size=UINT-64                   Batch size
    -s, --data-lenght=100                      Lenght of single item data(chars)
    -w, --write-ratio=0.5                      Write ratio (ex. 0.2 will result with 20% writes)
    -t, --timeout=1s                           Timeout for requests


    Note:
    If you don't provide the requests amount or duration limit program will continue running 
    indefinitely unless it is manually stopped by pressing `ctrl-c`. 


## What's next - TODO:
mvp:
- fix error handling, add summary of successful writes, reads (success, with error, nodata returned, timeout)
- more options
    - read preference (not only from primary) Primary, PrimaryPreferred, SecondaryPreferred, Secondary, Nearest
    - majority write
- more params with functionality:
    - cursor read
improvements:
- faster http client, https://github.com/valyala/fasthttp
- check if automaxprocs will give better performance, https://github.com/uber-go/automaxprocs
- ci to automically build package and dockerfile
- helmchart for easy multi instance load test and easier install

known issues:
- rate limit accuracy, current have 30ops deviation with bigger rps's
- deviation of write/read ration ~up to 3ops, better ratio distribution

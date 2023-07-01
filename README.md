# mongoload - workload driver for MongoDB 

## Introduction
The purpose of this tool is to simulate workloads to facilitate testing the failover capabilities of MongoDB clusters under load. This code, being an open-source project, is in its early development stage and likely contains various bugs.


## How to use:
1. Build image - `docker build -t mload .`
2. Run - `docker run mload -uri=http://localhost:21017 -req=10000`

### Current flags:
- `-uri` - mongo uri connection
- `-db` - database name (default: 'load_test')
- `-col` - collection name (default: 'load_test_coll')
- `-conn` - number of concurrent connections
- `-rps` - request per second limit
- `-d` - duration (ex. 10s, 1m) 
- `-req` - requests to perform (inserts by default)
- `-bs` - batch size (if set inserts will be in batches)
- `-dl` - length of item (default: 100, len of chars in single item to insert)
- `-wr` - write ratio (default: 0.5)

Note:
If you don't provide the operations amount(`-ops`) or duration(`-d`), the program will continue running indefinitely unless it is manually stopped by pressing `ctrl-c`. 


## What's next - TODO:
- add timeouts for queries
- fix error handling, add summary of successful writes, reads (success, with error, nodata returned, timeout)
- more options
    - read preference (not only from primary) Primary, PrimaryPreferred, SecondaryPreferred, Secondary, Nearest
    - majority write
- more params with functionality:
    - cursor read
- faster http client, https://github.com/valyala/fasthttp
- check if automaxprocs will give better performance, https://github.com/uber-go/automaxprocs
- ci to automically build package and dockerfile
- helmchart for easy multi instance load test and easier install

known issues:
- rate limit accuracy, current have 30ops deviation with bigger rps's
- deviation of write/read ration ~up to 3ops, better ratio distribution

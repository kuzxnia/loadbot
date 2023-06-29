# mongoload - workload driver for MongoDB 


The purpose of this tool is to simulate workloads and to facilitate test the failover capabilities of MongoDB clusters under the load. This code, being open source, is in its early stages and likely contains numerous bugs.

requirements:
* go in version 1.20 

current flags:

* `-uri` - mongo uri connection
* `-db` - database name (default: 'load_test')
* `-col` - collection name (default: 'load_test_coll')
* `-conn` - number of concurrent connections
* `-rps` - request per second limit
* `-d` - duration (ex. 10s, 1m) 
* `-req` - requests to perform (inserts by default)
* `-bs` - batch size (if set inserts will be in batches)
* `-dl` - length of item (default: 100, len of chars in single item to insert)

Note:
If you don't provide the operations amount(`-ops') or duration(`-d`), the program will run indefinitely.


todo:
* simpler build - makefile?
* more params with functionality:
    * ~~single write~~
    * batch write
    * single read
    * cursor read
    * mixed (if you provide will 50/50 by default)
        * mix ratio

* dockerfile
* ci to automically build package and dockerfile

* helmchart for multi instance load test and easier install

* ~~working rateLimit~~, added more accurate rate limit, current have 30ops deviation with bigger rps's

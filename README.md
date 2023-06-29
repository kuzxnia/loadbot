# mongoload - workload driver for MongoDB 


todo:
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

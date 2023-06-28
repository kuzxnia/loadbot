# mongoload

!!! NOT ready !!!

currently 

simple mongodb load provider tool

todo:
* working duration and indefinitely if not set (with and without ratelimit)
* dockerfile
* ci to automically build package and dockerfile
* more params with functionality:
    * single write
    * batch write
    * single read
    * cursor read
    * mixed (if you provide will 50/50 by default)
        * mix ratio

* helmchart for multi instance load test and easier install

* ~~working rateLimit~~, added more accurate rate limit, current have 30ops deviation with bigger rps's

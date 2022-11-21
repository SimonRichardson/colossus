# Request

The following request service aims to provide a way to manage requests over
http in a more distributed manor. The service takes one request and distributes
the same request over to multiple end points, using a replica setup.

It expects a primary service and then slave services can be defined during the
configuration.


                              +---------------+        +-----------+
                              |               |        |           |
                        +-----> CLUSTER (N)   +-------->  SERVICE  |
                        |     |               |        |           |
                        |     +---------------+        +-----------+
       +---------+      |     +---------------+        +-----------+
       |         |      |     |               |        |           |
       | REQUEST +------------> CLUSTER (N+1) +-------->  SERVICE  |
       |         |      |     |               |        |           |
       +---------+      |     +---------------+        +-----------+
                        |     +---------------+        +-----------+
                        |     |               |        |           |
                        +-----> CLUSTER (N+1) +-------->  SERVICE  |
                              |               |        |           |
                              +---------------+        +-----------+


The following provides a very light weight abstraction over resque with only
a very lightweight API being exposed (to date, see below):

 - Request

-----

## Examples

The request expects the url be the primary service url, all slave services will
be sent at the same time, but with a different url.

```go
req := http.NewRequest("GET", "http://unknown.com", nil)
resp, err := request.Request(res)
```

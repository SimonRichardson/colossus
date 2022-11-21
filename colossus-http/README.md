# Colossus http

------

The http server is a reasonable implementation for easy insertion of item ids
into the system.

------

## Usage

The server intended to be run as a standalone implementation.

### Environmental variables

The colossus application has a set of environment variables to help tweak the
application for different setups (testing vs production). The application can
use the various strategies (see root README.md) to then turn on an off various
parts of the application:

#### Local Development

When developing locally it's advised to run the services (redis, colossus, etc)
through docker (esp. through docker compose). To help with this (assuming you're
not using native docker), it's adviced to export the following environmental
variables locally to your shell.

```bash
export DOCKER_IP=$(docker-machine ip)
export MONGO_INSTANCES="$DOCKER_IP:27017"
export REDIS_INSTANCES="$DOCKER_IP:6377;$DOCKER_IP:6378;$DOCKER_IP:6379"
```

### Running

Running the colossus is relatively easy and can even be run side by side the
insert server by passing a different port to run on. If you just want to test
out the colossus application just run the following:

```bash
go run ./colossus-http/main.go
```

Alternatively running the application with a different port, then just overwrite
the environmental variable.

```bash
HTTP_ADDRESS=":9002" go run colossus-http/main.go
```

### API

Operations are differentiated by their HTTP verb. All endpoints must be sent
using flatbuffers protocol, to prevent slowdown with the std json marshalling
and unmarshalling.

Keys must be bson ObjectId hexs.

Note that write operations will claim success and return 200 as long as
the quorum is achieved, even if the provided score was lower than what has
already been persisted and therefore the operation was actually a `no-op`.

#### Select

GET to `/key/select?size=100&expiry=100`.

```bash
$ curl -XGET 'http://localhost:9001'
```

To read the response:
```go
s := &records.OKKeyFieldScoreTxnValue{}
s.Read(body)

for _, v := range s.Records {
    r := &records.PutRecord{}
    r.Read([]byte(v.Value[1:]))

    fmt.Println(m.OwnerId)
}
```

### Operations

Colossus expects to interact with a set of independent Redis instances, which
operators should deploy, monitor and manage. In general, Redis will use a lot of
RAM and comparatively little CPU and Colossus will use very little RAM and
comparatively large amount of CPU. It may make sense to co-locate a Colossus
instance with every Redis instance.

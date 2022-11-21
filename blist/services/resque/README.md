# Resque

The following resque service aims to provide a way to manage messages sent and
recieved over resque in a more distributed manor and by distributed I mean for
the CPU.

The following provides a very light weight abstraction over resque with only
a very lightweight API being exposed (to date, see below):

 - EnqueueBytes
 - DequeueBytes
 - RegisterFailure

-----

## Notes

When writing to redis it's performs a _write to all_ clusters and waits for
all clusters to respond before moving on. This setup allows for strong
consistency when _all clusters_ are available, but if there is a network
partition then the resque setup will become inconsistent. So with that in mind
only fire and forget messages should be sent over resque, if you want both
highly consistency along with availability then a repair strategy should be
employed to fix it.

Reading only reads from one cluster, so in a partition period then it could be
possible that messages could have been missed.

-----

## Examples

For enqueuing and dequeuing it's expected that you've already marshalled the
items at hand into a series of bytes so that it can enqueue it to the right
queue for the resque clients to handle.

```go
args := []string{"hello", "world", "!"}
bytes, _ := json.Marshal(args)

resque.EnqueueBytes(s.Queue("ping"), s.Class("Hello"), bytes)
```

When a failure happens, or you want to notify the message bus of a failure then
calling `RegisterFailure` is the best way to do so.

```go
failure := s.Failure{fmt.Errorf("Something went wrong!")}

resque.RegisterFailure(s.Queue("ping"), s.Class("Hello"), failure)
```

# Colossus

    ELI5: High-performance ephemeral ticket store (i.e. pre-DB tickets)

-----

This is *Colossus* a service for serving tickets at mass volume.

![](http://cdn.ttgtmedia.com/ITKE/cwblogs/public-sector/assets_c/2011/08/Colossus%20valves%20photographed%20by%20Andrew%20Back-thumb-500x357-138123.png)

-----

### Introduction

Colossus aims to coordinate between different services (mongodb, redis and
internal memory) to provide a reasonable ticketing service. The design goal from
the outset of Colossus is to be agnostic as possible about the payload of the
data.

Colossus implements different storage types, strategies depending on the
service it's utilizing. For example the counter service implements a time-series
event storage via a LWW-element-set CRDT with limited inline garbage collection.

At a high level the counter service maintains sets of values, with each set
ordered accordingly with an associated value timestamp. Similarly the store
service is implemented as a LWW-element-set, but does not fully embody a CRDT
completely because it doesn't fully abide by the following laws that the counter
does.

CRDTs (conflict-free replicated data types) are data types on which the same set
of operations yields the same outcome when performed regardless of order of the
executions and duplications of the operations. This allows data convergence
without the need for consensus between replicas. This allows for easier an
implementation, because no consensus protocol implementation is required.

Operations on CRDTs need to adhere to the following laws:

 - Associativity

   The grouping of operations don't matter.

   `a + (b + c) = (a + b) + c`

 - Commutativity

   The order of the operations don't matter.

   `a + b = b + a`

 - Idempotence

   Duplication of operations don't matter.

   `a + a = a`

Colossus implements a set data type, specifically the Last Writer Wins element
set (LWW-element-set). The description of the LWW-element-set simply follows:

 - An element is in the set, if its most-recent operation was an add.
 - An element is not in the set, if it's most-recent operation was a remove.


A more formal description of a LWW-element-set, as informed by [Shapiro](https://hal.inria.fr/file/index/docid/555588/filename/techreport.pdf),
is as follows:

 - Set `S` is represented by two internal sets, the add set `A` and the remove
 set `R`.
 - To add an element `e` to the set `S`, add a tuple `t` with the element and
 the current timestamp `t = (e, now())` to `A`.
 - To remove an element from the set `S`, add a tuple `t` with the element and
 the current timestamp `t = (e, now())` to `R`.
 - To check if an element `e` is in the set `S`, check if it is in the add set
 `A` and not in the remove set `R` with a higher timestamp.

Colossus implements the above definition, but extends it by applying a basic
asynchronous garbage collection. All nodes carry a timestamp which are then
filtered out of the nodes if the timestamp has become expired.

-----

### Servers

1. [HTTP Server](colossus-http/README.md)
1. [Walker Server](colossus-walker/README.md)

-----

### Replication

Colossus replicates data over multiple non-communicating clusters. A typical
replication factor is 3. Colossus has two methods of replicating data:

1. During a write.
1. During a read-repair.

A write (insert or delete) is sent to all clusters. The overall operation
returns success when the quorum is reached. Unsuccessful clusters might have
been affected by a network partition (slow, failed, crash) and in case of a
unsuccessful write then a read-repair might be triggered on a later read.

A read (select) is dependent on the read strategy employed. If the strategy
queries several clusters, it might be able to spot a disjointment in the
resulting sets. If so, the union set is returned to the client and in the
background, a read-repair is triggered which lazily converges the sets across
all the replicas.

-----

### Fault tolerance

Colossus runs as a homogenous distributed system. Each Colossus instance can
serve all requests (insert, delete, select) for a client, and communicates
with all Redis instances.

A Colossus instance is effectively stateless, but holds transient state. If a
Colossus instance crashes, three types of state are lost:

1. Current client connections are lost. Clients can reconnect to another
Colossus instance and re-execute their operations.
1. Lost client connections can lead to unfulfilled stored requests, where
different stores can hold inconsistent states. Rolling back and garbage
collection should clean these states up, but there isn't 100% guarantee at this
state.
1. Unresolved read-repairs are lost. The read repair might be triggered again
during another read.

Since all store operations are idempotent, failure modes should not impede on
convergence of the data.

Persistence is delegated to MongoDB (others to follow).

If a Redis instance is permanently lost and has to be replaces with a fresh
instance, there are two options:

1. Replace it with an empty instance. Keys will be replicated to it via a
read-repair. Aa more and more keys are replicated, the read-repair load will
decrease and the instance will work normally. This process might result in data
loss over the lifetime of a system. If the other replicas are also lost,
non-replicated keys (keys that have not been requested and thus did not trigger
a read-repair) are lost.
1. Replace it with a cloned replica. There will be a gap between the time of the
last write respected by the replica and the first write respected by the new
instance. The gap might be fixed by subsequent read-repairs.

Both processes can be expedited via a keyspace walker process. Nevertheless,
these properties and procedures warrant careful consideration.

-----

### Structure

The structure of Colossus sets out to be tunable depending on the work
undertaken (currently not possible without a restart). It is possible to change
the various strategies for each service so that a different approach can be
utilized (performance vs memory vs bandwidth).

1. Coordinator

  The coordinator's job is to work as a intermediary for all the various
services, similar in principale to conventional controller, but the difference
is that the coordinator has more of a role in managing the stores. If the
coordinator encounters an error it should either attempt to recover (by
repairing or rolling back) or fallback to another service when possible.

2. Farms

  A farm is a collection of clusters (see below) that allow the creation of fail
overs or improvement of speed with the aid of more servers. The Coordinator
generally talks to Farms directly.

3. Clusters

  A collection of services in a group, most of which are just pools of
connections. The cluster will do the main calls to the services (mongodb, redis
and etc) and can be tweaked to use normal serial commands or pipelined commands
to reduce the network latency.

4. Pools

  Pools hold a collection of connections directly to the service. To aid the
performance of each service, multiple connections are created to help improve
any latency issues whilst waiting for a request to come back, with the added
benefit of not exhausting the network stack of connections.

5. Instrumentation

  Through out the application as much instrumentation has been put in to the
application to ensure that we can see if any performance/defects occur. This can
be seen as either plaintext (files, stdout), statsd and etc.

-----

### Naming

1. [Colossus Computer](https://en.wikipedia.org/wiki/Colossus_computer)

  Colossus was the name of a series of computers developed by British
  codebreakers in 1943-1945 to help in the cryptanalysis of the Lorenz cipher.
  Colossus used thermionic valves (vacuum tubes) and thyratrons to perform
  Boolean and counting operations. Colossus is thus regarded[1] as the world's
  first programmable, electronic, digital computer, although it was programmed
  by plugs and switches and not by a stored program.

2. [I am Colossus](https://www.youtube.com/watch?v=fTYXbFsWg-M)

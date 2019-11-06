bookish-couscous
================

## Introduction

My own attempt to zen.ly hiring challenge (~2017/11/15).

The description is in ChallengeZenly.md.

Language used:

- golang

Software component used:

- kafka & zookeeper (wurstmeister/kafka, wurstmeister/zookeeper)
- scylladb (scylladb/scylla)
- redis (redis:alpine)
- nginx (as a tcp load balancer for frontends)

Libs used:

- https://github.com/golang/geo
- https://github.com/gocql/gocql
- https://github.com/bsm/sarama-cluster
- https://github.com/golang/protobuf
- https://godoc.org/google.golang.org/grpc
- https://github.com/Shopify/sarama
- https://github.com/garyburd/redigo/redis
- https://github.com/golang/protobuf/proto

Please note this is my first use of following technologies, so I mostly
followed basic tutorials and I'm maybe not using them perfectly:

- S2 geo library;
- scylla & gocql;
- kafka & shopify's sarama;
- protobuf
- grpc


## Up & running

```shell
# Build an initial image, with all deps, to make things faster
$ docker build -t bookish-couscous-base .

# Start the stack:
$ docker-compose up -d
[...]

$ docker-compose ps
            Name                           Command              State              Ports
----------------------------------------------------------------------------------------------------
bookishcouscous_client_1        /bin/sh -c /bin/sh /go/src      Up
                                ...
bookishcouscous_fo01_1          /bin/sh -c /go/bin/fo           Up
bookishcouscous_fo02_1          /bin/sh -c /go/bin/fo           Up
bookishcouscous_fo_1            /opt/nginx/sbin/nginx -g d      Up      443/tcp, 80/tcp
                                ...
bookishcouscous_generator_1     /bin/sh -c /bin/sh /go/src      Up
                                ...
bookishcouscous_kafka_1         start-kafka.sh                  Up      0.0.0.0:32794->9092/tcp
bookishcouscous_processor01_1   /go/bin/processor -init         Up
bookishcouscous_processor02_1   /bin/sh -c /go/bin/processor    Up
bookishcouscous_processor03_1   /bin/sh -c /go/bin/processor    Up
bookishcouscous_processor04_1   /bin/sh -c /go/bin/processor    Up
bookishcouscous_redis_1         docker-entrypoint.sh redis      Up      6379/tcp
                                ...
bookishcouscous_scylla_1        /docker-entrypoint.py           Up      10000/tcp, 7000/tcp,
                                                                        7001/tcp, 9042/tcp,
                                                                        9160/tcp, 9180/tcp
bookishcouscous_zookeeper_1     /bin/sh -c /usr/sbin/sshd       Up      0.0.0.0:2181->2181/tcp,
                                ...                                     22/tcp, 2888/tcp, 3888/tcp

# To send events:
$ docker-compose exec generator generator -h
Usage of generator:
  -days int
        Number of days (default 15)
  -event int
        Number of event to inject per day (default 5000)
  -friends int
        Number of friends per user (default 10)
  -kafka string
        kafka host port (default "kafka:9092")
  -redis string
        redis host port (default "redis:6379")
  -users int
        Number of users (default 100)

# Generate events...
$ docker-compose exec generator generator -days 7 -event 1000 -friends 5 -users 10
[...]

# Check in db...
$ docker-compose exec scylla cqlsh -e 'select * from zenly.kyf;'
[...]

# To query fo:
$ docker-compose exec client client -uid 42
2017/11/16 10:25:27 My uid is 42
2017/11/16 10:25:27 Best friend:          1
2017/11/16 10:25:27 Crush:                42
2017/11/16 10:25:27 Most seen:            1
2017/11/16 10:25:27 Mutual love:          1
2017/11/16 10:25:27 Mutual love all time: 1

# (note: if given uid = user id, then it is like it was not found;
# just a lazy thing from me :) )
```


## Components description

### generator

Generate friends, relationships, and sessions; It writes friends' relationships in redis,
and sessions in a kafka queue. It also generate & stores users' SP in redis (only one, home).
It will generate days * events number of events, over 100 users having each 10 friends (pick randomnly).

### processor

Fetch sessions from queue, and ran processing over it. It will then store in DB valuable information,
ready to be handled by fo on request.

Its main algo is:

for each session in queue:

- fetch session from queue;
- if session invalid (not friends), return;
- grab metadata from both user 1 & user 2 relationships with each others;
- compute time passed together on session, store it in metadata struct according of session's day ("most seen" feature);
- grab both users SP from redis;
- if relevant, store time passed together for "best friend" feature;
- store time passed for "most seen, all time" feature
- if night & position requirements are fullfilled, store night in structure as well ("crush feature");
- save both structure in database;
- done!

Note that "mutual love" feature is not computed here. It'll be computed when information will
be necessary by FO (The reason is that when the FO call occurs, this information may be
already outdated by relationships, and thus will be needed to be recomputed...)

It runs only 4 db queries: 2 read & 2 writes (1 for each users of the session).
It also runs 2 read on redis.
It will invalidate also data computed & stored by fo component (see next section).

It can be ran multiple instance of processors. Each are independent and keeps no persistant information
in memory.

### fo

Listens for client's requests through a grpc socket.

On a client request, it will fetch user's friends data from scylla and compute through a single iteration
over its friend's relationships metadata for best friends, most seens, etc.

At the end of the process, we'll do the same over its "most seen" (7 days & all time) friend,
and if uids are the same on each sides of relationship, they we can conclude "mutual love".

Result is stored in redis (cached 12h - because daily data become obsolete at a point) & returned by grpc socket.

### client

Opens a grpc socket, queries the fo.


## DB Schema

I've used only 1 scylla table for this:

```sql
CREATE TABLE IF NOT EXISTS zenly.kyf (
    user_id bigint,
    rel_user_id bigint,
    PRIMARY KEY(user_id, rel_user_id),
    duration bigint,
    week_most bigint,
    week_friends bigint,
    nights list<timestamp>,
    week_most_list map<timestamp, int>,
    week_friends_list map<timestamp, int>,
);
```

Note that the primary key is over (userid, reluserid). I do not know if this is
efficient enough to do scan queries over all userid records (that I do on the
fo part). I'm pretty sure this is something to investigate, perform some
perfomance tests & read docs about it.

Sample queries ran over this schema:

```sql
SELECT rel_user_id, week_friends_list, week_most_list, nights, duration
    FROM kyf
    WHERE user_id = ?
```

```
SELECT duration, nights, week_most_list, week_friends_list
    FROM kyf WHERE user_id = ? AND rel_user_id = ?
```

Data sample:

```sql
cqlsh> select * from zenly.kyf;

 user_id | rel_user_id | duration | nights                       | week_friends_list                                                                                                                                                                                                                              | week_most_list
---------+-------------+----------+------------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
       0 |           1 |    15440 | ['2017-11-14 06:00:35+0000'] | {'2017-11-11 00:00:00+0000': 1973, '2017-11-12 00:00:00+0000': 2127, '2017-11-13 00:00:00+0000': 1569, '2017-11-14 00:00:00+0000': 1818, '2017-11-15 00:00:00+0000': 2516, '2017-11-16 00:00:00+0000': 2348, '2017-11-17 00:00:00+0000': 1131} | {'2017-11-11 00:00:00+0000': 2167, '2017-11-12 00:00:00+0000': 2127, '2017-11-13 00:00:00+0000': 1569, '2017-11-14 00:00:00+0000': 2289, '2017-11-15 00:00:00+0000': 2516, '2017-11-16 00:00:00+0000': 2348, '2017-11-17 00:00:00+0000': 1236}
       1 |           0 |    15440 | ['2017-11-14 06:00:35+0000'] |  {'2017-11-11 00:00:00+0000': 2167, '2017-11-12 00:00:00+0000': 2127, '2017-11-13 00:00:00+0000': 1569, '2017-11-14 00:00:00+0000': 1913, '2017-11-15 00:00:00+0000': 1916, '2017-11-16 00:00:00+0000': 2065, '2017-11-17 00:00:00+0000': 788} | {'2017-11-11 00:00:00+0000': 2167, '2017-11-12 00:00:00+0000': 2127, '2017-11-13 00:00:00+0000': 1569, '2017-11-14 00:00:00+0000': 2289, '2017-11-15 00:00:00+0000': 2516, '2017-11-16 00:00:00+0000': 2348, '2017-11-17 00:00:00+0000': 1236}
```

## Known drawbacks & possible enhancements

- timezones are not managed;
- long sessions are not managed either (multiple nights);
- data generator could be enhanced to simulate day periods in days & nights, create similar
  patterns, and make a better use of SPs
- May require better integration testing
- database records locking issue (2 processors working on same users will lead to data overwrite.)
  A solution would be not to write data in scylla from processors, but send computations done by
  processors in a kafka topic & have a last unique process to store everything in cassandra. We
  would not store values, but add them to stored (update kyf set duration = duration + nn where uid = 42)
  Another would be to lock row records, but scylla doesn't allow this. Doing it elsewhere seems dangerous.

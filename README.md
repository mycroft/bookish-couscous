bookish-couscous
================

## Introduction

My own attempt to zen.ly hiring challenge (~2017/11/15).

- Challenge: https://gist.github.com/daedric/db45c531a1bc5e58f0383f9c1bff4306
- Hosted on https://gitlab.mkz.me/mycroft/bookish-couscous

Language used:

- golang

Software component used:

- kafka & zookeeper (wurstmeister/kafka, wurstmeister/zookeeper)
- scylladb (scylladb/scylla)
- redis (redis:alpine)

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

# To send events:
$ docker exec -ti bookishcouscous_generator_1 generator -h
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

# To query fo:
$ docker exec -ti bookishcouscous_client_1 client 42
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

## Known drawbacks & possible enhancements

- timezones are not managed;
- long sessions are not managed either (multiple days/nights);
- data generator could be enhanced to simulate day periods in days & nights, create similar
  patterns, and make a better use of SPs
- there is a lot of things that could be cached on FO side
- multiple fo!!!
- May require better integration testing


## Existing tests

- TBW


## Todo:

- Revoir l'algo nuit
- Data generator
- multiple fo (google tcp backend), multiple processor

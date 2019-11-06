# Know Your Friends

Being a social network, Zenly is to provide to its users useful information 
about their friends. Of course this has to be derived from what we receive.

One valuable information is to categorize the friends in several categories
which can later be reused or simply displayed on the client (the app).

Categorizing the friends depends on a very specific dimension, we call
_significant places_. The data science provides a data store with the Home,
School/Work of our users.

Depending on the time of the day and where you spent some time with your friends,
we are able to further improve our sorting.

## Your project

*You are to get all of our users' sessions and categorize friend according to
those sessions.*

### Environment

You can use the language you prefer amongst: Go, C++, Scala, Java or Python.

#### Data available
Today, we have a stream of _sessions_, sessions represents the time two friends
spent together in a specific location. In this stream you can expect:

* user id 1;
* user id 2;
* starting date;
* end date;
* latitude;
* longitude;

Those events are published in a Kafka topic using a Protobuf schema.

There is also a Redis which:

* Given a user id will return his significant places;
* Given a user id will return all his friends.

#### Databases

The only distributed data store available is ScyllaDB.

## Requirements

Given those categories:

* Best Friend: the person you see the most outside of *your* significant place;
* Crush: Person with whom you have spent 3 nights in the 7 days at either their
  home or yours, have different homes. The night is specified as following:
  from 22h to 8am and session of at least 6 hours;
* Most seen: the person you see the most;
* Mutual Love: You've spent most time with them and they have spent most time
with you.


We want the result for all of those categories for the last _rolling_ 7 days.
We want the Mutual Love also calculated on all the sessions we know.


## Expected solution

* Specify all the protobuf schemas you need (hint: store protobuf in Redis);
* Write a quick and dirty session producer;
* Have one service that can be queried from the client (thus real time, use
  gRPC) that gives back the category list filled with the friend id;
* In case of restart services are always expected to be in a consistent state
  upon restart;
* Of course, you are not to write only one service. Feel free to write 
  services/component you need.

Good luck.

## Resources

* You should most probably look at S2:
  
  * http://blog.christianperone.com/2015/08/googles-s2-geometry-on-the-sphere-cells-and-hilbert-curve/
  * https://github.com/golang/geo
  * Should you decide to use S2, cells of level 16 are precise enough.
  
* As Kafka producer/consumer we use: https://github.com/bsm/sarama-cluster
* To interact with ScyllaDB we use: https://github.com/gocql/gocql


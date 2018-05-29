# Redis Proxy for Segment

A Redis Proxy service implemented as a HTTP web service.

## Architecture / Code Overview

The proxy utilizes the Go net/http package to listen for incomming GET requests.

For the LRU cache I decided to use karlseguin's ccache, which already has support for concurrency. I made the assumption that using a library for this was alright. See the following section for more discussion.

To connect to the redis I use the radix.v2 library to create a pool of redis connections and then get and return one of the connections for every outgoing command I make to the redis.

## Algorighmic Complexity of Cache Operations

Here I will describe the algorithmic complexity for operations using the cache in various situations.

A very basic LRU will have a map to keep track of keys and map to values, which will be stored in a linked list to keep track of ordering. CCache adds the idea of buckets, which shard it's internal map to provide a greater ammount of concurrency. 

### Cache Hit O(1)

In the case of a cache hit, aka when the item we need is already in the cache and not expired, we get the correct bucket from by accessing an index of an array based on a lookup from a map, then get the value from the map that is stored in the bucket. Although in the worst case getting a value from a map can be O(n), the usual case and how I've usually done this in the past we simplify this to 0(1) runtime.

CCache uses a buffered channel to queue promotions to a single worker, in order to reduce lock contention. So the item is now queued for promotion, which also takes constant time.

When the time comes to do promtions, the worker thread moves the relevent items to the front of the list, which I assume in Go takes constant time.

### Cache Miss O(m)

In the case of a cache miss, aka when the item we need is not in the cache or is expired, the first thing we need to do is fetch the item from the redis client. We will say this takes O(m) time.

CCcache then creates a new entry with the provided expiration and inserts it into the map of the approprate bucket, as well as adding it into the list which tracks recency, which takes constant time.

If the size has been reached, we go from the back of the list and delete least recently used entries from the map and the list. Since it's a doubly linked list, this takes constant time.

So the overall algorithmic complexity O(m) if you are accounting for the required read from redis, and O(1) otherwise.

## Instructions for Running tests

Step 1. Fix the docker stuff

Step 2. Run `make test`

Alternatively:
Step 1. Start a redis with address localhost:6379 and start main.go also on localhost with port 8080.

Step 2. Run `make test`

## What I spent time on

I have not previously used Go, worked with Redis, set up a docker container, or used docker-compose.

I spent about an hour or two making things work locally, including writing the proxy server and e2e tests.

I spent what feels like 12 hours on all the docker build things.

## What I didn't implment

The system tests for testing concurrency and that the LRU is an LRU and that it is following the key size constraint could be more thourough. I'm not sure that an end-to-end test is where I would usually put these things... Also I don't have any more time to allocate to this.

The platform requirement and the single click build and test requirement because I can't allocate any more time to this.



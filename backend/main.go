// Generate binary data from static web assets
//go:generate go-bindata ../app/dist/...

package main

import (
	"flag"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

// Session for AWS SDK calls
//var Session *session.Session

func main() {
	profile := flag.String("profile", "default", "a string")
	region := flag.String("region", "us-east-1", "a string")
	flag.Parse()

	c := ufo.UFOConfig {
		Profile: profile,
		Region: region,
	}

	u := ufo.Fly(c)

	router := routes(u)
	router.Run()
}

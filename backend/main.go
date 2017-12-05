// Generate binary data from static web assets
//go:generate go-bindata ../app/dist/...

package main

import (
	"flag"
	"fmt"

	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

type logger struct{}

func (l *logger) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func main() {
	profile := flag.String("profile", "default", "a string")
	region := flag.String("region", "us-east-1", "a string")
	flag.Parse()

	c := ufo.UFOConfig{
		Profile: profile,
		Region:  region,
	}

	u := ufo.Fly(c, &logger{})

	router := routes(u)
	router.Run()
}

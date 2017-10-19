package main

import (
	"flag"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
	"github.com/abiosoft/ishell"
	log "github.com/sirupsen/logrus"
)

func main() {
	profile := flag.String("profile", "", "a string")
	region := flag.String("region", "us-east-1", "a string")

	flag.Parse()

	if *profile == "" {
		log.Fatalln("Profile option required.")
	}

	c := ufo.UFOConfig {
		Profile: profile,
		Region: region,
	}

	app := &App{
		UFO: ufo.Fly(c),
		Shell: ishell.New(),
	}

	app.Init()
}

package main

import (
	"flag"
	"gitlab.fuzzhq.com/Web-Ops/ufo/ufo"
	"github.com/abiosoft/ishell"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	profile := flag.String("profile", "", "AWS profile to use")
	region := flag.String("region", "us-east-1", "AWS region to use")

	noInteractive := flag.Bool("i", false, "Non-interactive deployment.")

	cluster := flag.String("cluster", "", "Cluster to deploy to. Required if non-interactive.")
	service := flag.String("service", "", "Service to deploy to. Required if non-interactive.")
	version := flag.String("version", "", "Version to deploy. Required if non-interactive.")

	flag.Parse()

	if *profile == "" {
		log.Fatalln("Profile option required.")
	}

	if *noInteractive {
		for _, val := range []*string{cluster, service, version} {
			if *val == "" {
				log.Printf("Missing required option.")
				flag.PrintDefaults()
				os.Exit(1)
			}
		}
	}

	c := ufo.UFOConfig {
		Profile: profile,
		Region: region,
	}

	app := &App{
		UFO: ufo.Fly(c),
		Shell: ishell.New(),
		f: &AppFlags{
			*noInteractive,
			*cluster,
			*service,
			*version,
		},
	}

	app.Init()
}

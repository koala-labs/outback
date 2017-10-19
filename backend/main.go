package main

import (
	"flag"

	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/ecr"
	//"github.com/aws/aws-sdk-go/service/ecs"
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

	//Session = session.Must(session.NewSessionWithOptions(session.Options{
	//	Config:  aws.Config{Region: aws.String(*region)},
	//	Profile: *profile,
	//}))
	//ECSService = ecs.New(Session)
	//ECRService = ecr.New(Session)

	router := routes(u)
	router.Run()
}

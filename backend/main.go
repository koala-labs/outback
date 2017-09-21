package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// Session for AWS SDK calls
var Session = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

func main() {
	router := routes()
	router.Run()
}

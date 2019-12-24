package main

import (
	"github.com/kimmyfek/next_rtd/database/dynamo"

	"github.com/aws/aws-lambda-go/lambda"
)

var al *dynamo.AccessLayer

func init() {
	var err error
	al, err = dynamo.New()
	if err != nil {
		// TODO return err
	}
}

func main() {
	lambda.Start(al.GetStationsAndConnections)
}

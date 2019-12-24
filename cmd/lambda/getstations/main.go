package main

import (
	"encoding/json"
	"net/http"

	"github.com/kimmyfek/next_rtd/database/dynamo"

	"github.com/aws/aws-lambda-go/events"
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

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	s, err := al.GetStationsAndConnections()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	j, err := json.Marshal(s)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(j),
	}, nil
}

func main() {
	lambda.Start(show)
}

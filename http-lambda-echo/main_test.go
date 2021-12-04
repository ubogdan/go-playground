package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo/v4"
)

func TestLambdaHandler(t *testing.T) {
	r := echo.New()
	r.GET("/", testHandler)
	lambda := LambdaHandler(r)

	res, err := lambda(context.Background(), events.APIGatewayProxyRequest{
		Resource:                        "",
		Path:                            "/",
		HTTPMethod:                      "GET",
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  nil,
		StageVariables:                  nil,
		RequestContext:                  events.APIGatewayProxyRequestContext{},
		Body:                            "",
		IsBase64Encoded:                 false,
	})
	if err != nil {
		t.Error("Failed %w", err)
	}
	t.Logf("response %#v", res)
}

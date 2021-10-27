package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
)

func testHandler(c echo.Context) error {
	return json.NewEncoder(c.Response().Writer).Encode(map[string]string{
		"Message": "This is a JSON message",
	})
}

func LambdaHandler(r *echo.Echo) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Proxy handler
	proxy := echoadapter.New(r)

	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return proxy.ProxyWithContext(ctx, request)
	}
}

func main() {

	r := echo.New()
	r.GET("/", testHandler)

	// Start lambda event handler
	lambda.Start(LambdaHandler(r))
}

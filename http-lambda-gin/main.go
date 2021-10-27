package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

func testHandler(ctx *gin.Context) {
	ctx.JSON(200, map[string]string{
		"Message": "This is a JSON message",
	})
}

func LambdaHandler(r *gin.Engine) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Proxy handler
	proxy := ginadapter.New(r)

	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return proxy.ProxyWithContext(ctx, request)
	}
}

func main() {
	r := gin.New()
	r.GET("/", testHandler)

	// Start lambda event handler
	lambda.Start(LambdaHandler(r))
}

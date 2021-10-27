package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
)

func testHandler(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{
		"Message": "This is a JSON message",
	})
}

func LambdaHandler(r *mux.Router) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Proxy handler
	proxy := gorillamux.New(r)

	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return proxy.ProxyWithContext(ctx, request)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", testHandler)

	// Start lambda event handler
	lambda.Start(LambdaHandler(r))
}

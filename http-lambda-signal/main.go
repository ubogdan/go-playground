package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
)

func testHandler(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{
		"Message": "This is a JSON message",
	})
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", testHandler)

	lambda.StartWithOptions(gorillamux.New(r), lambda.WithEnableSIGTERM(func() {
		log.Printf("Closing database connection")
	}))
}

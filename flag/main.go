package main

import (
	"flag"
	"log"
	"time"
)

var (
	username, password string
	maxRetry           int
	randomUserAgent    bool
	timeout            time.Duration
)

func main() {
	flag.IntVar(&maxRetry, "maxRetry", 3, "Number of retries")
	flag.StringVar(&username, "u", "", "Basic authetication user")
	flag.StringVar(&password, "p", "", "Basic authentication password")
	flag.DurationVar(&timeout, "t", 1*time.Second, "sleep period")

	flag.BoolVar(&randomUserAgent, "ua", false, "")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("url is required param !")
	}

	// Test with: ./flag -maxRetry 10 -ua -t 5s https://www.google.ro
	log.Printf("retrieve url %s for max %d with randomAgent=%v and a timeout of %s", args[0], maxRetry, randomUserAgent, timeout)

	// Test with: ./flag -u john -p doe https://api.github.com/v3
	if len(username) > 0 && len(password) > 0 {
		log.Printf("with basic auth credentils %s:%s", username, password)
	}
}

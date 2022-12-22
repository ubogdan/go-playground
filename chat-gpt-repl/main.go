package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
)

const maxTokens = 1024

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing API KEY")
	}

	client := gpt3.NewClient(apiKey)

	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("Me> ")
	for reader.Scan() {
		resp, err := client.CompletionWithEngine(context.Background(),
			"text-davinci-002",
			gpt3.CompletionRequest{
				Prompt:    []string{reader.Text()},
				MaxTokens: gpt3.IntPtr(maxTokens),
			})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("GPT3>", resp.Choices[0].Text)
		fmt.Print("Me> ")
	}
	fmt.Println()

}

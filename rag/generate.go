package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

func Generate(query string, chunks []string) (string, error) {
	prompt := fmt.Sprintf(template, query, strings.Join(chunks, "\n\n"))
	fmt.Printf("\nprompts: ==============================\n%s", prompt)
	return ollamaGenerate(prompt)
}

func ollamaGenerate(prompt string) (string, error) {
	var response string
	stream := false
	err := ollamaClient.Generate(context.Background(), &api.GenerateRequest{
		Model:  "qwen3.5:cloud",
		Prompt: prompt,
		Stream: &stream,
		Think:  &api.ThinkValue{Value: false},
	}, func(gr api.GenerateResponse) error {
		response = gr.Response
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("ollama generate response error: %w", err)
	}
	return response, nil
}

package rag

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

func EmbeddingChunks(chunks []string) ([][]float64, error) {
	vectors := make([][]float64, 0, len(chunks))
	for _, chunk := range chunks {
		vector, err := ollamaEmbddding(chunk)
		if err != nil {
			return nil, err
		}
		vectors = append(vectors, vector)
	}

	return vectors, nil
}

func EmbeddingChunk(chunk string) ([]float64, error) {
	vector, err := ollamaEmbddding(chunk)
	if err != nil {
		return nil, err
	}
	return vector, nil
}

func ollamaEmbddding(chunk string) ([]float64, error) {
	resp, err := ollamaClient.Embeddings(context.Background(), &api.EmbeddingRequest{
		Model:  "qwen3-embedding",
		Prompt: chunk,
	})
	if err != nil {
		return nil, fmt.Errorf("embedding chunk error: %w", err)
	}
	return resp.Embedding, nil
}

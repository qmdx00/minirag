package rag

import (
	"context"
	"fmt"
	"runtime"
	"strconv"

	chromem "github.com/philippgille/chromem-go"
)

var (
	collectionName = "knowledge-base"
	DB             *chromem.DB
)

func init() {
	if DB == nil {
		DB = chromem.NewDB()
	}
}

func SaveToDB(chunks []string, vectors [][]float64) error {
	collection, err := DB.CreateCollection(collectionName, nil, nil)
	if err != nil {
		return fmt.Errorf("create db collection error: %w", err)
	}

	ctx := context.Background()
	documents := make([]chromem.Document, 0, len(chunks))

	for index := range chunks {
		vec64 := vectors[index]
		vec32 := make([]float32, 0, len(vec64))
		for _, v := range vec64 {
			vec32 = append(vec32, float32(v))
		}

		documents = append(documents, chromem.Document{
			ID:        strconv.Itoa(index),
			Embedding: vec32,
			Content:   chunks[index],
		})
	}

	if err := collection.AddDocuments(ctx, documents, runtime.NumCPU()); err != nil {
		return fmt.Errorf("add collection documents error: %w", err)
	}

	return nil
}

func QueryFromDB(queryEmbedding []float64, topK int) ([]chromem.Result, error) {
	vec32 := make([]float32, 0, len(queryEmbedding))
	for _, v := range queryEmbedding {
		vec32 = append(vec32, float32(v))
	}

	collection := DB.GetCollection(collectionName, nil)
	results, err := collection.QueryEmbedding(context.Background(), vec32, topK, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("query documents error: %w", err)
	}
	return results, nil
}

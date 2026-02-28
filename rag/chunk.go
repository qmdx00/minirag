package rag

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func SplitIntoChunks(docPath string) []string {
	content, err := getFileContent(docPath)
	if err != nil {
		return nil
	}

	chunks := make([]string, 0)
	for chunk := range strings.SplitSeq(content, "\n\n") {
		if chunk == "" {
			continue
		}
		chunks = append(chunks, chunk)
	}

	return chunks
}

func getFileContent(docPath string) (string, error) {
	file, err := os.Open(docPath)
	if err != nil {
		return "", fmt.Errorf("open file failed, err: %w, path= %s", err, docPath)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("read file content error: %w", err)
	}

	return string(content), nil
}

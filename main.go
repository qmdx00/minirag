package main

import (
	"fmt"
	"sort"

	"github.com/philippgille/chromem-go"
	"github.com/qmdx00/minirag/rag"
)

// RAG (Retrieval-Augmented Generation 检索增强生成) 步骤：
//
// 提问前准备:
//  1. 分片: 通过一定的方式对内容进行分片（如按照段落，章节等方式切分）
//  2. 索引：向量化+存储，通过 embedding 模型对分片内容进行向量化，并存储到向量数据库中。
//
// 提问后处理:
//  1. 召回：根据用户提示词从向量数据库召回 topN 条数据，根据相似度排序（成本低，速度快，准确度低）
//  2. 重排：使用 cross encoder 模型对召回后结果进行重排，选出 topK 个最相关内容(成本高，耗时长，准确度高)
//  3. 生成：将重排后的内容与用户提示词一起输入到大模型，获取输出结果。

func main() {
	// 分片
	chunks := rag.SplitIntoChunks("./docs/milvus.txt")

	// 索引: 使用 embedding 模型对分片文本进行向量化
	vectors, err := rag.EmbeddingChunks(chunks)
	if err != nil {
		panic(err)
	}

	// 存储到向量数据库
	if err := rag.SaveToDB(chunks, vectors); err != nil {
		panic(err)
	}

	// 用户提问，提示词向量化
	question := `什么是向量数据库`
	vector, err := rag.EmbeddingChunk(question)
	if err != nil {
		panic(err)
	}

	// 召回 top 10
	results, err := rag.QueryFromDB(vector, 10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nTOP5 results: ====================\n")
	printResults(results)

	// 根据相似度排序召回结果
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Cross encoder 模型重排 (暂不进行重排，直接取召回相似度高的3条)
	fmt.Printf("\nTOP3 results: ====================\n")
	rerankResults := results[:3]
	printResults(rerankResults)

	// 整合内容输入到大模型,获取生成结果
	rerankChunks := make([]string, 0, len(rerankResults))
	for _, result := range rerankResults {
		rerankChunks = append(rerankChunks, result.Content)
	}
	response, err := rag.Generate(question, rerankChunks)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nGenerate response: =====================\n")
	fmt.Println(response)
}

func printResults(results []chromem.Result) {
	for _, result := range results {
		fmt.Printf("ID=%s, Similarity=%f, Content=%s\n", result.ID, result.Similarity, result.Content)
	}
}

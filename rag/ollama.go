package rag

import "github.com/ollama/ollama/api"

var ollamaClient *api.Client

func init() {
	if ollamaClient == nil {
		ollamaClient, _ = api.ClientFromEnvironment()
	}
}

var template = `你是一位知识助手，请根据用户的问题和下列片段生成准确的回答。
用户问题：%s
相关片段：
%s
请基于上述内容作答，不要编造信息。"！！"`

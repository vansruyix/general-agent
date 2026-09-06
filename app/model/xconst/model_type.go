package xconst

import "strings"

// ModelType 模型类型枚举
type ModelType string

// modelTypeDefinition 定义模型类型及其对应的原始类型
type modelTypeDefinition struct {
	modelType  ModelType
	originType string
}

// 所有支持的模型类型定义
var modelTypeDefinitions = []modelTypeDefinition{
	{LLM, "text-generation"},
	{TextEmbedding, "embeddings"},
	{Rerank, "reranking"},
	{Speech2Text, "speech2text"},
	{TTS, "tts"},
	{Text2Img, "text2img"},
	{Moderation, "moderation"},
}

const (
	LLM           ModelType = "llm"            // 系统推理模型
	TextEmbedding ModelType = "text-embedding" // 知识库文档嵌入处理的默认模型
	Rerank        ModelType = "rerank"         // 重排序模型
	Speech2Text   ModelType = "speech2text"    // 语言转文本模型
	TTS           ModelType = "tts"            // 文本转语音模型
	Text2Img      ModelType = "text2img"       // 文本转图片模型
	Moderation    ModelType = "moderation"     // 文本审核模型
)

// ToOriginModelType 将ModelType转换为原始模型类型字符串
func (m ModelType) ToOriginModelType() string {
	for _, def := range modelTypeDefinitions {
		if def.modelType == m {
			return def.originType
		}
	}
	return ""
}

// ToModelType 将原始模型类型字符串转换为ModelType
func ToModelType(originType string) ModelType {
	for _, def := range modelTypeDefinitions {
		if def.originType == originType || string(def.modelType) == originType {
			return def.modelType
		}
	}
	return ""
}

// GetModelType 获取模型类型
func GetModelType(modelType string) ModelType {
	for _, def := range modelTypeDefinitions {
		if strings.Contains(modelType, string(def.modelType)) || strings.Contains(modelType, def.originType) {
			return def.modelType
		}
	}
	return ""
}

package model

// RAGflowKnowledge RAGflow知识库
type KnowledgeBasicInfo struct {
	Cancelled  int `json:"cancelled"`
	Downloaded int `json:"downloaded"`
	Failed     int `json:"failed"`
	Finished   int `json:"finished"`
	Processing int `json:"processing"`
}

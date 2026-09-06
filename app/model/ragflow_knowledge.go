package model

import "general-agent/extension/xtime"

// RAGflowKnowledge RAGflow知识库
type RAGflowKnowledge struct {
	Id          string         `json:"id" gorm:"primarykey"` // 知识库ID
	UpdatedDate xtime.JsonTime `json:"updated_at"`           // 更新时间
	TenantId    string         `json:"tenant_id"`            // 租户ID
	Nickname    string         `json:"nickname"`             // 创建者昵称
	Name        string         `json:"name"`                 // 知识库名称
	Language    string         `json:"language"`             // 语言
	Description string         `json:"description"`          // 描述
	EmbdId      string         `json:"embd_id"`              // 嵌入模型ID
	ParserId    string         `json:"parser_id"`            // 解析器ID
	Permission  string         `json:"permission"`           // 权限
	DocNum      int            `json:"doc_num"`              // 文档数
	TokenNum    int            `json:"token_num"`            // 总token数
	ChunkNum    int            `json:"chunk_num"`            // 总分块数
	UpdateTime  int64          `json:"update_time"`          // 更新时间戳
	Avatar      *string        `json:"avatar"`               // 头像（可为空）
}

func (RAGflowKnowledge) TableName() string {
	return "knowledgebase"
}

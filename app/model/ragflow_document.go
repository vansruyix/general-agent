package model

import (
	"encoding/json"
	"general-agent/extension/xtime"
)

// RAGflowDocument RAGflow知识库文档
type RAGflowDocument struct {
	Id            string         `json:"id" gorm:"primarykey"` // 文档ID
	CreateTime    int64          `json:"create_time"`          // 创建时间戳
	CreateTimeStr string         `json:"create_date"`          // 创建时间字符串
	UpdateTime    int64          `json:"update_time"`          // 更新时间戳
	UpdateDate    xtime.JsonTime `json:"update_date"`          // 更新时间
	Thumbnail     string         `json:"thumbnail"`            // 缩略图
	KbId          string         `json:"kb_id"`                // 知识库ID
	ParserId      string         `json:"parser_id"`            // 解析器ID
	// PipelineId      string          `json:"pipeline_id"`          // 流水线ID
	ParserConfig    json.RawMessage `json:"parser_config" swaggertype:"object"` // 解析器配置（JSON格式）
	SourceType      string          `json:"source_type"`                        // 文档来源类型
	Type            string          `json:"type"`                               // 文档类型
	CreatedBy       string          `json:"created_by"`                         // 创建者
	Nickname        string          `json:"nickname"`                           // 创建者昵称
	Name            string          `json:"name"`                               // 文档名称
	Location        string          `json:"location"`                           // 位置
	Size            int             `json:"size"`                               // 文档大小，单位B
	TokenNum        int             `json:"token_num"`                          // token数量
	ChunkNum        int             `json:"chunk_num"`                          // 分块数
	Progress        float64         `json:"progress"`                           // 进度
	ProgressMsg     string          `json:"progress_msg"`                       // 进度信息
	ProcessBeginAt  xtime.JsonTime  `json:"process_begin_at"`                   // 开始时间
	ProcessDuration float64         `json:"process_duration"`                   // 处理时长
	MetaFields      json.RawMessage `json:"meta_fields" swaggertype:"object"`   // 元数据（JSON格式）
	Suffix          string          `json:"suffix"`                             // 文件后缀
	Run             string          `json:"run"`                                // 任务状态
	Status          string          `json:"status"`                             // 文档状态
}

func (RAGflowDocument) TableName() string {
	return "document"
}

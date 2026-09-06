package xconst

import "net/http"

// DifyEndpoints dify功能接口对照表
var DifyEndpoints = map[Func]*Endpoint{
	Login:     NewEndpoint(http.MethodPost, "/login"),
	Apps:      NewEndpoint(http.MethodPost, "/apps"),
	GetApps:   NewEndpoint(http.MethodGet, "/apps"),
	UpdateApp: NewEndpoint(http.MethodPut, "/apps/{appId}"),
	DeleteApp: NewEndpoint(http.MethodDelete, "/apps/{appId}"),

	MemberList:     NewEndpoint(http.MethodGet, "/workspaces/current/members"),                  // 成员列表
	InviteEmail:    NewEndpoint(http.MethodPost, "/workspaces/current/members/invite-email"),    // 邀请邮箱
	UpdateRole:     NewEndpoint(http.MethodPut, "/workspaces/current/members/{id}/update-role"), // 更新角色
	DeleteUser:     NewEndpoint(http.MethodDelete, "/workspaces/current/members/{id}"),          // 删除成员
	UpdatePassword: NewEndpoint(http.MethodPost, "/account/password"),                           // 更新密码

	ListModelByType:              NewEndpoint(http.MethodGet, "/workspaces/current/models/model-types/{model_type}"),                   // 根据模型类型获取模型列表
	ListModelByProvider:          NewEndpoint(http.MethodGet, "/workspaces/current/model-providers/{provider}/models"),                 // 根据模型供应商获取模型列表
	AddModelByProvider:           NewEndpoint(http.MethodPost, "/workspaces/current/model-providers/{provider}/models"),                // 根据模型供应商新增模型
	GetProviderCredential:        NewEndpoint(http.MethodGet, "/workspaces/current/model-providers/{provider}/credentials"),            // 根据模型供应商获取凭据
	GetModelCredentialByProvider: NewEndpoint(http.MethodGet, "/workspaces/current/model-providers/{provider}/models/credentials"),     // 根据模型供应商获取模型凭证
	GetModelParamRuleByProvider:  NewEndpoint(http.MethodGet, "/workspaces/current/model-providers/{provider}/models/parameter-rules"), // 根据模型供应商获取模型参数规则
	EnableModelByProvider:        NewEndpoint(http.MethodPatch, "/workspaces/current/model-providers/{provider}/models/enable"),        // 启用模型
	DisableModelByProvider:       NewEndpoint(http.MethodPatch, "/workspaces/current/model-providers/{provider}/models/disable"),       // 禁用模型
	DeleteModelByProvider:        NewEndpoint(http.MethodDelete, "/workspaces/current/model-providers/{provider}/models"),              // 删除模型
	SetProviderApiKey:            NewEndpoint(http.MethodPost, "/workspaces/current/model-providers/{provider}"),                       // 设置模型供应商API密钥
	SetDefaultModel:              NewEndpoint(http.MethodPost, "/workspaces/current/default-model"),

	PublishWorkflowApp: NewEndpoint(http.MethodPost, "/apps/{id}/workflows/publish"), // 发布应用

	ListTools: NewEndpoint(http.MethodGet, "/workspaces/current/tool-providers"), // 获取工具列表

	UpdateToolProviderCredential: NewEndpoint(http.MethodPost, "/workspaces/current/tool-provider/builtin/{provider}/update"),

	DatasetsExternal:       NewEndpoint(http.MethodPost, "/datasets/external"),
	DatasetsExternalApi:    NewEndpoint(http.MethodPost, "/datasets/external-knowledge-api"),
	DatasetsExternalDelete: NewEndpoint(http.MethodDelete, "/datasets/{id}"),

	GetWorkflowDraft:         NewEndpoint(http.MethodGet, "/apps/{appId}/workflows/draft"),
	SyncWorkflowDraft:        NewEndpoint(http.MethodPost, "/apps/{appId}/workflows/draft"),
	RunNode:                  NewEndpoint(http.MethodPost, "/apps/{appId}/workflows/draft/nodes/{nodeId}/run"),
	RunIterationNode:         NewEndpoint(http.MethodPost, "/apps/{appId}/workflows/draft/iteration/nodes/{nodeId}/run"),
	RunLoopNode:              NewEndpoint(http.MethodPost, "/apps/{appId}/workflows/draft/loop/nodes/{nodeId}/run"),
	RunAdvancedChatWorkflow:  NewEndpoint(http.MethodPost, "/apps/{appId}/advanced-chat/workflows/draft/run"),
	RunWorkflow:              NewEndpoint(http.MethodPost, "/apps/{appId}/workflows/draft/run"),
	GetNodeDefaultConfig:     NewEndpoint(http.MethodGet, "/apps/{appId}/workflows/default-workflow-block-configs/{blockType}"),
	GetAllNodeDefaultConfigs: NewEndpoint(http.MethodGet, "/apps/{appId}/workflows/default-workflow-block-configs"),
	GetWorkflowConfig:        NewEndpoint(http.MethodGet, "/apps/{appId}/workflows/draft/config"),
	GetPublishedWorkflow:     NewEndpoint(http.MethodGet, "/apps/{appId}/workflows/publish"),
	PublishWorkflow:          NewEndpoint(http.MethodPost, "/apps/{appId}/workflows/publish"),
	GetWorkflowVersions:      NewEndpoint(http.MethodGet, "/apps/{appId}/workflows"),
	UpdateWorkflowVersion:    NewEndpoint(http.MethodPatch, "/apps/{appId}/workflows/{workflowId}"),
	DeleteWorkflowVersion:    NewEndpoint(http.MethodDelete, "/apps/{appId}/workflows/{workflowId}"),
	ImportWorkflowDSL:        NewEndpoint(http.MethodPost, "/apps/{appId}/workflows/draft/import"),
	StopWorkflowRun:          NewEndpoint(http.MethodPost, "/apps/{appId}/workflows/runs/{taskId}/stop"),
	GetWorkflowRunDetail:     NewEndpoint(http.MethodGet, "/apps/{appId}/workflow-runs/{runId}"),
	GetNodeExecutions:        NewEndpoint(http.MethodGet, "/apps/{appId}/workflow-runs/{runId}/node-executions"),
	GetConversationVariables: NewEndpoint(http.MethodGet, "/apps/{appId}/workflows/conversation-variables"),
}

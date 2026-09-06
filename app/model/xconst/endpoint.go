package xconst

type Func string

const (
	Register  Func = "register"
	Login     Func = "login"
	Apps      Func = "apps"
	GetApps   Func = "get_apps"
	UpdateApp Func = "update_app"
	DeleteApp Func = "delete_app"

	InviteEmail    Func = "create_user"
	MemberList     Func = "get_members"
	UpdateRole     Func = "update_role"
	DeleteUser     Func = "delete_user"
	UpdatePassword Func = "update_password"

	ListModelByType              Func = "list_model_by_type"
	ListModelByProvider          Func = "list_model_by_provider"
	AddModelByProvider           Func = "add_model_by_provider"
	GetProviderCredential        Func = "get_provider_credential"
	GetModelCredentialByProvider Func = "get_model_credential_by_provider"
	GetModelParamRuleByProvider  Func = "get_model_param_rule_by_provider"
	EnableModelByProvider        Func = "enable_model_by_provider"
	DisableModelByProvider       Func = "disable_model_by_provider"
	DeleteModelByProvider        Func = "delete_model_by_provider"
	SetProviderApiKey            Func = "set_provider_api_key"
	SetDefaultModel              Func = "set_default_model"

	PublishWorkflowApp Func = "publish_app"

	ListTools Func = "list_tools"

	UpdateToolProviderCredential Func = "update_tool_provider_credential"

	GetLlmList             Func = "get_llm_list"
	GetLlmFactories        Func = "get_llm_factories"
	GetMyLlms              Func = "get_my_llms"
	SetLlmApiKey           Func = "set_llm_api_key"
	AddLlm                 Func = "add_llm"
	DeleteLlm              Func = "delete_llm"
	EnableLlm              Func = "enable_llm"
	DeleteLlmFactory       Func = "delete_llm_factory"
	SetTenantInfo          Func = "set_tenant_info"
	GetTenantInfo          Func = "get_tenant_info"
	AddKnowledge           Func = "add_knowledge"
	UpdateKnowledge        Func = "update_knowledge"
	DeleteKnowledge        Func = "delete_knowledge"
	ListKnowledge          Func = "list_knowledge"
	DetailKnowledge        Func = "detail_knowledge"
	BasicKnowledge         Func = "basic_knowledge"
	ListPipelineLog        Func = "list_pipeline_log"
	ListPipelineDatasetLog Func = "list_pipeline_dataset_log"
	ListTags               Func = "list_tags"

	AddDocument           Func = "add_document"
	DownloadDocument      Func = "download_document"
	UploadDocument        Func = "upload_document"
	RenameDocument        Func = "rename_document"
	ChangeDocumentStatus  Func = "change_document_status"
	RunDocument           Func = "run_document"
	DeleteDocument        Func = "delete_document"
	ChangePaeser          Func = "change_parser"
	SetDocumentMetaFeilds Func = "set_document_meta_fields"

	FindChunk      Func = "find_chunk"
	CreateChunk    Func = "create_chunk"
	SwitchChunk    Func = "switch_chunk"
	RetrievalChunk Func = "retrieval_chunk"
	DeleteChunk    Func = "delete_chunk"
	DetailChunk    Func = "detail_chunk"
	UpdateChunk    Func = "update_chunk"
	GetChunkImg    Func = "get_chunk_img"

	CreateToken Func = "create_token"

	DatasetsExternal       Func = "datasets_external"
	DatasetsExternalApi    Func = "datasets_external_api"
	DatasetsExternalDelete Func = "datasets_external_delete"

	GetWorkflowDraft         Func = "get_workflow_draft"
	SyncWorkflowDraft        Func = "sync_workflow_draft"
	RunNode                  Func = "run_node"
	RunIterationNode         Func = "run_iteration_node"
	RunLoopNode              Func = "run_loop_node"
	RunAdvancedChatWorkflow  Func = "run_advanced_chat_workflow"
	RunWorkflow              Func = "run_workflow"
	GetNodeDefaultConfig     Func = "get_node_default_config"
	GetAllNodeDefaultConfigs Func = "get_all_node_default_configs"
	GetWorkflowConfig        Func = "get_workflow_config"
	GetPublishedWorkflow     Func = "get_published_workflow"
	PublishWorkflow          Func = "publish_workflow"
	GetWorkflowVersions      Func = "get_workflow_versions"
	UpdateWorkflowVersion    Func = "update_workflow_version"
	DeleteWorkflowVersion    Func = "delete_workflow_version"
	ImportWorkflowDSL        Func = "import_workflow_dsl"
	StopWorkflowRun          Func = "stop_workflow_run"
	GetWorkflowRunDetail     Func = "get_workflow_run_detail"
	GetNodeExecutions        Func = "get_node_executions"
	GetConversationVariables Func = "get_conversation_variables"
)

const (
	GetVerifyAuthUrl = "http://dimrealm-vsm-verify-auth-file:5539/api/v1/verify-auth/authorize-new" // 获取授权信息
)

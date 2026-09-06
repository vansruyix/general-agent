/*
 * @Author       : dufei
 * @Date         : 2025-12-18 17:44:01
 * @LastEditors  : dufei
 * @LastEditTime : 2026-01-09 15:02:47
 * @Description  :
 */

package xconst

var RAGflowEndpoints = map[Func]*Endpoint{
	Register: NewEndpoint("POST", "/user/register"),
	Login:    NewEndpoint("POST", "/user/login"),

	GetLlmList:             NewEndpoint("GET", "/llm/list"),
	GetLlmFactories:        NewEndpoint("GET", "/llm/factories"),
	GetMyLlms:              NewEndpoint("GET", "/llm/my_llms"),
	SetLlmApiKey:           NewEndpoint("POST", "/llm/set_api_key"),
	AddLlm:                 NewEndpoint("POST", "/llm/add_llm"),
	DeleteLlm:              NewEndpoint("POST", "/llm/delete_llm"),
	EnableLlm:              NewEndpoint("POST", "/llm/enable_llm"),
	DeleteLlmFactory:       NewEndpoint("POST", "/llm/delete_factory"),
	SetTenantInfo:          NewEndpoint("POST", "/user/set_tenant_info"),
	GetTenantInfo:          NewEndpoint("GET", "/user/tenant_info"),
	AddKnowledge:           NewEndpoint("POST", "/datasets"),
	UpdateKnowledge:        NewEndpoint("PUT", "/datasets/{kb_id}"),
	DeleteKnowledge:        NewEndpoint("DELETE", "/datasets"),
	ListKnowledge:          NewEndpoint("GET", "/datasets?"),
	DetailKnowledge:        NewEndpoint("GET", "/kb/detail?kb_id={kb_id}"),
	BasicKnowledge:         NewEndpoint("GET", "/kb/basic_info?kb_id={kb_id}"),
	ListPipelineLog:        NewEndpoint("POST", "/kb/list_pipeline_logs?kb_id={kb_id}&page={page}&page_size={page_size}&keywords={keywords}"),
	ListPipelineDatasetLog: NewEndpoint("POST", "/kb/list_pipeline_dataset_logs?kb_id={kb_id}&page={page}&page_size={page_size}&keywords={keywords}"),
	ListTags:               NewEndpoint("GET", "/kb/tags"),

	AddDocument:           NewEndpoint("POST", "/document/create"),
	DownloadDocument:      NewEndpoint("GET", "/document/get/{doc_id}"),
	UploadDocument:        NewEndpoint("POST", "/document/upload"),
	RenameDocument:        NewEndpoint("POST", "/document/rename"),
	ChangeDocumentStatus:  NewEndpoint("POST", "/document/change_status"),
	RunDocument:           NewEndpoint("POST", "/document/run"),
	DeleteDocument:        NewEndpoint("POST", "/document/rm"),
	ChangePaeser:          NewEndpoint("POST", "/document/change_parser"),
	SetDocumentMetaFeilds: NewEndpoint("POST", "/document/set_meta"),

	FindChunk:      NewEndpoint("POST", "/chunk/list"),
	CreateChunk:    NewEndpoint("POST", "/chunk/create"),
	SwitchChunk:    NewEndpoint("POST", "/chunk/switch"),
	RetrievalChunk: NewEndpoint("POST", "/chunk/retrieval_test"),
	DeleteChunk:    NewEndpoint("POST", "/chunk/rm"),
	DetailChunk:    NewEndpoint("GET", "/chunk/get"),
	UpdateChunk:    NewEndpoint("POST", "/chunk/set"),
	GetChunkImg:    NewEndpoint("GET", "/document/image/{img_id}"),

	CreateToken: NewEndpoint("POST", "/system/new_token"),
}

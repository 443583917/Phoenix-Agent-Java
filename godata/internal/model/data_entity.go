package model

// ===== Data Module GORM Entities =====
//
// Each entity embeds BaseModel for common fields:
// ID (string PK), CreateTime, UpdateTime, CreateBy, UpdateBy, DelFlag.
// Table names match the Java MyBatis Flex @Table annotations exactly.
// Field types and JSON names match the Java entity sources.
//
// Java enum fields (KnowledgeType, EmbeddingStatus, ModelType, AgentStatusEnm,
// SessionStatusEnm) are mapped to Go string types.
//
// NOTE: Some Java String fields named "tableName" are renamed to "Table" in Go
// to avoid conflict with GORM's Tabler interface (TableName() method).

// ----------------------------------------------------------------
// 1. Agent — tbl_data_agent
// ----------------------------------------------------------------

type Agent struct {
	BaseModel
	Type          string `gorm:"column:type;type:varchar(32);default:sql" json:"type"`
	CategoryId    string `gorm:"column:category_id;type:varchar(64)" json:"categoryId"`
	Sn            string `gorm:"column:sn;type:varchar(64)" json:"sn"`
	Name          string `gorm:"column:name;type:varchar(128)" json:"name"`
	Description   string `gorm:"column:description;type:text" json:"description"`
	Avatar        string `gorm:"column:avatar;type:varchar(256)" json:"avatar"`
	Status        string `gorm:"column:status;type:varchar(32)" json:"status"`
	ApiKey        string `gorm:"column:api_key;type:varchar(256)" json:"-"`
	ApiKeyEnabled int    `gorm:"column:api_key_enabled;default:0" json:"apiKeyEnabled"`
	Prompt        string `gorm:"column:prompt;type:text" json:"prompt"`
	Category      string `gorm:"column:category;type:varchar(64)" json:"category"`
	AdminId       int64  `gorm:"column:admin_id" json:"adminId"`
	Tags          string `gorm:"column:tags;type:varchar(512)" json:"tags"`
	OrderNum      int    `gorm:"column:order_num;default:0" json:"orderNum"`
}

func (Agent) TableName() string { return "tbl_data_agent" }

// ----------------------------------------------------------------
// 2. AgentCategory — tbl_data_agent_category
// Extends Java com.phoenix.common.model.BaseModel (creator, updator, delFlag).
// In Go, these map to BaseModel.CreateBy, BaseModel.UpdateBy, BaseModel.DelFlag.
// ----------------------------------------------------------------

type AgentCategory struct {
	BaseModel
	Pid         string `gorm:"column:pid;type:varchar(64)" json:"pid"`
	Sn          string `gorm:"column:sn;type:varchar(64)" json:"sn"`
	Name        string `gorm:"column:name;type:varchar(128)" json:"name"`
	Description string `gorm:"column:description;type:varchar(256)" json:"description"`
}

func (AgentCategory) TableName() string { return "tbl_data_agent_category" }

// ----------------------------------------------------------------
// 3. AgentDatasource — tbl_data_agent_datasource
// Many-to-many link between Agent and Datasource.
// ----------------------------------------------------------------

type AgentDatasource struct {
	BaseModel
	AgentId      int64 `gorm:"column:agent_id" json:"agentId"`
	DatasourceId int   `gorm:"column:datasource_id" json:"datasourceId"`
	IsActive     int   `gorm:"column:is_active;default:1" json:"isActive"`
}

func (AgentDatasource) TableName() string { return "tbl_data_agent_datasource" }

// ----------------------------------------------------------------
// 4. AgentDatasourceTables — tbl_data_agent_datasource_tables
// Specific tables selected from a datasource for an agent.
// ----------------------------------------------------------------

type AgentDatasourceTables struct {
	BaseModel
	AgentDatasourceId int    `gorm:"column:agent_datasource_id" json:"agentDatasourceId"`
	Table             string `gorm:"column:table_name;type:varchar(256)" json:"tableName"`
}

func (AgentDatasourceTables) TableName() string { return "tbl_data_agent_datasource_tables" }

// ----------------------------------------------------------------
// 5. AgentKnowledge — tbl_data_agent_knowledge
// Agent knowledge documents, QA pairs, and FAQ entries.
// ----------------------------------------------------------------

type AgentKnowledge struct {
	BaseModel
	AgentId           int    `gorm:"column:agent_id" json:"agentId"`
	Title             string `gorm:"column:title;type:varchar(256)" json:"title"`
	Type              string `gorm:"column:type;type:varchar(32)" json:"type"`
	Question          string `gorm:"column:question;type:text" json:"question"`
	Content           string `gorm:"column:content;type:text" json:"content"`
	IsRecall          int    `gorm:"column:is_recall;default:0" json:"isRecall"`
	EmbeddingStatus   string `gorm:"column:embedding_status;type:varchar(32)" json:"embeddingStatus"`
	ErrorMsg          string `gorm:"column:error_msg;type:text" json:"errorMsg"`
	SourceFilename    string `gorm:"column:source_filename;type:varchar(256)" json:"sourceFilename"`
	FilePath          string `gorm:"column:file_path;type:varchar(512)" json:"filePath"`
	FileSize          int64  `gorm:"column:file_size;default:0" json:"fileSize"`
	FileType          string `gorm:"column:file_type;type:varchar(64)" json:"fileType"`
	SplitterType      string `gorm:"column:splitter_type;type:varchar(32);default:token" json:"splitterType"`
	IsResourceCleaned int    `gorm:"column:is_resource_cleaned;default:0" json:"isResourceCleaned"`
}

func (AgentKnowledge) TableName() string { return "tbl_data_agent_knowledge" }

// ----------------------------------------------------------------
// 6. AgentPresetQuestion — tbl_data_agent_preset_question
// ----------------------------------------------------------------

type AgentPresetQuestion struct {
	BaseModel
	AgentId   int64  `gorm:"column:agent_id" json:"agentId"`
	AccountId string `gorm:"column:account_id;type:varchar(64)" json:"accountId"`
	Question  string `gorm:"column:question;type:text" json:"question"`
	Answer    string `gorm:"-" json:"answer,omitempty"`
	SortOrder int    `gorm:"column:sort_order;default:0" json:"sortOrder"`
	IsActive  *bool  `gorm:"column:is_active;default:false" json:"isActive"`
}

func (AgentPresetQuestion) TableName() string { return "tbl_data_agent_preset_question" }

// ----------------------------------------------------------------
// 7. BusinessKnowledge — tbl_data_business_knowledge
// Business term definitions for LLM context.
// ----------------------------------------------------------------

type BusinessKnowledge struct {
	BaseModel
	BusinessTerm    string `gorm:"column:business_term;type:varchar(256)" json:"businessTerm"`
	Description     string `gorm:"column:description;type:text" json:"description"`
	Synonyms        string `gorm:"column:synonyms;type:varchar(512)" json:"synonyms"`
	IsRecall        int    `gorm:"column:is_recall;default:1" json:"isRecall"`
	AgentId         int64  `gorm:"column:agent_id" json:"agentId"`
	EmbeddingStatus string `gorm:"column:embedding_status;type:varchar(32)" json:"embeddingStatus"`
	ErrorMsg        string `gorm:"column:error_msg;type:text" json:"errorMsg"`
}

func (BusinessKnowledge) TableName() string { return "tbl_data_business_knowledge" }

// ----------------------------------------------------------------
// 8. ChatMessage — tbl_data_chat_message
// ----------------------------------------------------------------

type ChatMessage struct {
	BaseModel
	SessionId   string `gorm:"column:session_id;type:varchar(64)" json:"sessionId"`
	Role        string `gorm:"column:role;type:varchar(32)" json:"role"`
	Content     string `gorm:"column:content;type:text" json:"content"`
	MessageType string `gorm:"column:message_type;type:varchar(32)" json:"messageType"`
	Metadata    string `gorm:"column:metadata;type:text" json:"metadata"`
}

func (ChatMessage) TableName() string { return "tbl_data_chat_message" }

// ----------------------------------------------------------------
// 9. ChatSession — tbl_data_chat_session
// ----------------------------------------------------------------

type ChatSession struct {
	BaseModel
	AgentId  int    `gorm:"column:agent_id" json:"agentId"`
	Title    string `gorm:"column:title;type:varchar(256)" json:"title"`
	Status   string `gorm:"column:status;type:varchar(32);default:active" json:"status"`
	IsPinned bool   `gorm:"column:is_pinned;default:false" json:"isPinned"`
	UserId   string `gorm:"column:user_id;type:varchar(64)" json:"userId"`
}

func (ChatSession) TableName() string { return "tbl_data_chat_session" }

// ----------------------------------------------------------------
// 10. Datasource — tbl_data_datasource
// ----------------------------------------------------------------

type Datasource struct {
	BaseModel
	Name          string `gorm:"column:name;type:varchar(128)" json:"name"`
	Type          string `gorm:"column:type;type:varchar(32)" json:"type"`
	Host          string `gorm:"column:host;type:varchar(256)" json:"host"`
	Port          int    `gorm:"column:port" json:"port"`
	DatabaseName  string `gorm:"column:database_name;type:varchar(128)" json:"databaseName"`
	Username      string `gorm:"column:username;type:varchar(128)" json:"username"`
	Password      string `gorm:"column:password;type:varchar(256)" json:"-"`
	ConnectionUrl string `gorm:"column:connection_url;type:varchar(512)" json:"-"`
	Status        string `gorm:"column:status;type:varchar(32)" json:"status"`
	TestStatus    string `gorm:"column:test_status;type:varchar(32)" json:"testStatus"`
	Description   string `gorm:"column:description;type:varchar(256)" json:"description"`
	CreatorId     int64  `gorm:"column:creator_id" json:"creatorId"`
}

func (Datasource) TableName() string { return "tbl_data_datasource" }

// ----------------------------------------------------------------
// 11. LogicalRelation — tbl_data_logical_relation
// Logical foreign key configuration helping LLM understand data relationships.
// ----------------------------------------------------------------

type LogicalRelation struct {
	BaseModel
	DatasourceId     int    `gorm:"column:datasource_id" json:"datasourceId"`
	SourceTableName  string `gorm:"column:source_table_name;type:varchar(128)" json:"sourceTableName"`
	SourceColumnName string `gorm:"column:source_column_name;type:varchar(128)" json:"sourceColumnName"`
	TargetTableName  string `gorm:"column:target_table_name;type:varchar(128)" json:"targetTableName"`
	TargetColumnName string `gorm:"column:target_column_name;type:varchar(128)" json:"targetColumnName"`
	RelationType     string `gorm:"column:relation_type;type:varchar(16)" json:"relationType"`
	Description      string `gorm:"column:description;type:varchar(512)" json:"description"`
}

func (LogicalRelation) TableName() string { return "tbl_data_logical_relation" }

// ----------------------------------------------------------------
// 12. ModelConfig — tbl_data_model_config
// LLM model configuration for providers, endpoints, and proxy settings.
// ----------------------------------------------------------------

type ModelConfig struct {
	BaseModel
	Provider        string  `gorm:"column:provider;type:varchar(64)" json:"provider"`
	BaseUrl         string  `gorm:"column:base_url;type:varchar(512)" json:"baseUrl"`
	ApiKey          string  `gorm:"column:api_key;type:varchar(512)" json:"-"`
	ModelName       string  `gorm:"column:model_name;type:varchar(128)" json:"modelName"`
	Temperature     float64 `gorm:"column:temperature;type:decimal(3,2);default:0" json:"temperature"`
	IsActive        bool    `gorm:"column:is_active;default:false" json:"isActive"`
	MaxTokens       int     `gorm:"column:max_tokens;default:0" json:"maxTokens"`
	ModelType       string  `gorm:"column:model_type;type:varchar(32)" json:"modelType"`
	CompletionsPath string  `gorm:"column:completions_path;type:varchar(256)" json:"completionsPath"`
	EmbeddingsPath  string  `gorm:"column:embeddings_path;type:varchar(256)" json:"embeddingsPath"`
	ProxyEnabled    *bool   `gorm:"column:proxy_enabled;default:false" json:"proxyEnabled"`
	ProxyHost       string  `gorm:"column:proxy_host;type:varchar(256)" json:"proxyHost"`
	ProxyPort       int     `gorm:"column:proxy_port" json:"proxyPort"`
	ProxyUsername   string  `gorm:"column:proxy_username;type:varchar(128)" json:"proxyUsername"`
	ProxyPassword   string  `gorm:"column:proxy_password;type:varchar(256)" json:"-"`
}

func (ModelConfig) TableName() string { return "tbl_data_model_config" }

// ----------------------------------------------------------------
// 13. SemanticModel — tbl_data_semantic_model
// Semantic model mapping physical DB fields to business concepts.
// ----------------------------------------------------------------

type SemanticModel struct {
	BaseModel
	AgentId             int64  `gorm:"column:agent_id" json:"agentId"`
	DatasourceId        int    `gorm:"column:datasource_id" json:"datasourceId"`
	Table               string `gorm:"column:table_name;type:varchar(128)" json:"tableName"`
	ColumnName          string `gorm:"column:column_name;type:varchar(128)" json:"columnName"`
	BusinessName        string `gorm:"column:business_name;type:varchar(256)" json:"businessName"`
	Synonyms            string `gorm:"column:synonyms;type:varchar(512)" json:"synonyms"`
	BusinessDescription string `gorm:"column:business_description;type:text" json:"businessDescription"`
	ColumnComment       string `gorm:"column:column_comment;type:varchar(512)" json:"columnComment"`
	DataType            string `gorm:"column:data_type;type:varchar(64)" json:"dataType"`
	Status              int    `gorm:"column:status;default:1" json:"status"`
}

func (SemanticModel) TableName() string { return "tbl_data_semantic_model" }

// ----------------------------------------------------------------
// 14. UserPromptConfig — tbl_data_user_prompt_config
// User-defined prompt template configurations per agent or globally.
// ----------------------------------------------------------------

type UserPromptConfig struct {
	BaseModel
	Name         string `gorm:"column:name;type:varchar(128)" json:"name"`
	PromptType   string `gorm:"column:prompt_type;type:varchar(64)" json:"promptType"`
	AgentId      int64  `gorm:"column:agent_id" json:"agentId"`
	SystemPrompt string `gorm:"column:system_prompt;type:text" json:"systemPrompt"`
	Enabled      bool   `gorm:"column:enabled;default:true" json:"enabled"`
	Description  string `gorm:"column:description;type:varchar(256)" json:"description"`
	Priority     int    `gorm:"column:priority;default:0" json:"priority"`
	DisplayOrder int    `gorm:"column:display_order;default:0" json:"displayOrder"`
}

func (UserPromptConfig) TableName() string { return "tbl_data_user_prompt_config" }

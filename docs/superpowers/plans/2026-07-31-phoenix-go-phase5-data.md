# Phase 5 — phoenix-data (NL2SQL Engine + Agent Management)

**Scope:** 15 Java Controllers, 14 Entities. Split into 4 sub-phases:

### 5A: Data Entities
14 GORM entities: Agent, AgentCategory, AgentDatasource, AgentDatasourceTables, AgentKnowledge, AgentPresetQuestion, BusinessKnowledge, ChatMessage, ChatSession, Datasource, LogicalRelation, ModelConfig, SemanticModel, UserPromptConfig

### 5B: Agent Management CRUD (5 controllers)
AgentController, AgentCategoryController, AgentDatasourceController, AgentKnowledgeController, AgentPresetQuestionController — all standard CRUD patterns

### 5C: Chat + Session + Graph (3 controllers)
ChatController (SSE), GraphController (SSE NL2SQL), SessionEventController — streaming + session management

### 5D: Config CRUD + NL2SQL Workflows
DatasourceController, ModelConfigController, PromptConfigController, SemanticModelController, BusinessKnowledgeController, FileUploadController, EchoController — config CRUD + file upload

StateGraph workflow nodes in `agent/workflows/nl2sql/nodes/` — intent, evidence, schema, planner, sql_gen, python_exec, report

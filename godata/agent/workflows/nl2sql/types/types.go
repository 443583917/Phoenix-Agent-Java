package types

import "context"

// ──────────────────────────── State Key Constants ────────────────────────────

const (
	StateKeyInput                          = "input"
	StateKeyAgentID                        = "agent_id"
	StateKeyMultiTurnContext               = "multi_turn_context"
	StateKeyIntentOutput                   = "intent_recognition_node_output"
	StateKeyEvidence                       = "evidence"
	StateKeyQueryEnhanceOutput             = "query_enhance_node_output"
	StateKeyTableDocumentsForSchema        = "table_documents_for_schema_output"
	StateKeyColumnDocumentsForSchema       = "column_documents_for_schema_output"
	StateKeyTableRelationOutput            = "table_relation_output"
	StateKeyTableRelationException         = "table_relation_exception_output"
	StateKeyTableRelationRetryCount        = "table_relation_retry_count"
	StateKeyDBDialectType                  = "db_dialect_type"
	StateKeyGeneratedSemanticModelPrompt   = "generated_semantic_model_prompt"
	StateKeyFeasibilityOutput              = "feasibility_assessment_node_output"
	StateKeyPlannerOutput                  = "planner_node_output"
	StateKeyPlanCurrentStep                = "plan_current_step"
	StateKeyPlanNextNode                   = "plan_next_node"
	StateKeyPlanValidationStatus           = "plan_validation_status"
	StateKeyPlanValidationError            = "plan_validation_error"
	StateKeyPlanRepairCount                = "plan_repair_count"
	StateKeySQLGenerateOutput              = "sql_generate_output"
	StateKeySQLGenerateCount               = "sql_generate_count"
	StateKeySQLRegenReason                 = "sql_regenerate_reason"
	StateKeySemanticConsistencyOutput      = "semantic_consistency_node_output"
	StateKeySQLExecuteOutput               = "sql_execute_node_output"
	StateKeySQLResultListMemory            = "sql_result_list_memory"
	StateKeyPythonGenerateOutput           = "python_generate_node_output"
	StateKeyPythonExecuteOutput            = "python_execute_node_output"
	StateKeyPythonAnalyzeOutput            = "python_analysis_node_output"
	StateKeyPythonIsSuccess                = "python_is_success"
	StateKeyPythonTriesCount               = "python_tries_count"
	StateKeyPythonFallbackMode             = "python_fallback_mode"
	StateKeyIsOnlyNL2SQL                   = "is_only_nl2sql"
	StateKeyHumanReviewEnabled             = "human_review_enabled"
	StateKeyHumanFeedbackData              = "human_feedback_data"
	StateKeyResult                         = "result"
)

// ──────────────────────────── DTOs ────────────────────────────

// IntentRecognitionOutput is the output of the intent recognition node.
type IntentRecognitionOutput struct {
	Classification string `json:"classification"` // "闲聊或无关指令" or "可能的数据分析请求"
	Confidence     float64 `json:"confidence"`
	Reasoning      string  `json:"reasoning,omitempty"`
}

// QueryEnhanceOutput is the output of the query enhancement node.
type QueryEnhanceOutput struct {
	CanonicalQuery  string   `json:"canonical_query"`
	ExpandedQueries []string `json:"expanded_queries,omitempty"`
}

// Plan is the output of the planner node.
type Plan struct {
	ThoughtProcess string          `json:"thought_process"`
	ExecutionPlan  []ExecutionStep `json:"execution_plan"`
}

// ExecutionStep is a single step in the execution plan.
type ExecutionStep struct {
	Step           int               `json:"step"`
	ToolToUse      string            `json:"tool_to_use"`
	ToolParameters map[string]string `json:"tool_parameters"`
}

// SqlRetryDto holds SQL retry context.
type SqlRetryDto struct {
	Reason         string `json:"reason"`
	SemanticFail   bool   `json:"semantic_fail"`
	SQLExecuteFail bool   `json:"sql_execute_fail"`
}

// ──────────────────────────── Graph State ────────────────────────────

// NL2SQLState holds all state flowing through the NL2SQL graph nodes.
// This mirrors the Java OverAllState keys.
type NL2SQLState struct {
	// Input
	Input            string `json:"input"`
	AgentID          string `json:"agent_id"`
	MultiTurnContext string `json:"multi_turn_context"`

	// Intent
	IntentOutput *IntentRecognitionOutput `json:"intent_output,omitempty"`

	// Evidence
	Evidence string `json:"evidence,omitempty"`

	// Query enhancement
	QueryEnhanceOutput *QueryEnhanceOutput `json:"query_enhance_output,omitempty"`

	// Schema / Table relation
	TableDocumentsForSchema  []string `json:"table_documents_for_schema,omitempty"`
	ColumnDocumentsForSchema []string `json:"column_documents_for_schema,omitempty"`
	TableRelationOutput      string   `json:"table_relation_output,omitempty"`
	TableRelationException   string   `json:"table_relation_exception,omitempty"`
	TableRelationRetryCount  int      `json:"table_relation_retry_count"`

	// DB
	DBDialectType                string `json:"db_dialect_type,omitempty"`
	GeneratedSemanticModelPrompt string `json:"generated_semantic_model_prompt,omitempty"`

	// Feasibility
	FeasibilityOutput string `json:"feasibility_output,omitempty"`

	// Plan
	PlannerOutput         *Plan  `json:"planner_output,omitempty"`
	PlanCurrentStep       int    `json:"plan_current_step"`
	PlanNextNode          string `json:"plan_next_node,omitempty"`
	PlanValidationStatus  string `json:"plan_validation_status,omitempty"`
	PlanValidationError   string `json:"plan_validation_error,omitempty"`
	PlanRepairCount       int    `json:"plan_repair_count"`

	// SQL
	SQLGenerateOutput         string       `json:"sql_generate_output,omitempty"`
	SQLGenerateCount          int          `json:"sql_generate_count"`
	SQLRegenReason            *SqlRetryDto `json:"sql_regen_reason,omitempty"`
	SemanticConsistencyOutput string       `json:"semantic_consistency_output,omitempty"`
	SQLExecuteOutput          string       `json:"sql_execute_output,omitempty"`
	SQLResultListMemory       []string     `json:"sql_result_list_memory,omitempty"`

	// Python
	PythonGenerateOutput string `json:"python_generate_output,omitempty"`
	PythonExecuteOutput  string `json:"python_execute_output,omitempty"`
	PythonAnalyzeOutput  string `json:"python_analyze_output,omitempty"`
	PythonIsSuccess      bool   `json:"python_is_success"`
	PythonTriesCount     int    `json:"python_tries_count"`
	PythonFallbackMode   string `json:"python_fallback_mode,omitempty"`

	// Mode
	IsOnlyNL2SQL       bool   `json:"is_only_nl2sql"`

	// Human review
	HumanReviewEnabled bool   `json:"human_review_enabled"`
	HumanFeedbackData  string `json:"human_feedback_data,omitempty"`

	// Result
	Result string `json:"result,omitempty"`

	// Runtime metadata
	CurrentNode string   `json:"current_node,omitempty"`
	Errors      []string `json:"errors,omitempty"`
}

// ──────────────────────────── Legacy Types (for backward compat) ────────────────────────────

// EvidenceContext represents a piece of recalled evidence.
type EvidenceContext struct {
	Content string  `json:"content"`
	Score   float64 `json:"score"`
	Source  string  `json:"source"`
}

// SchemaContext represents a recalled schema element.
type SchemaContext struct {
	TableName    string `json:"tableName"`
	ColumnName   string `json:"columnName"`
	DataType     string `json:"dataType"`
	BusinessName string `json:"businessName"`
	Description  string `json:"description"`
}

// ExecutionPlan is the legacy plan type for backward compatibility.
type ExecutionPlan struct {
	Steps     []PlanStep `json:"steps"`
	Reasoning string     `json:"reasoning"`
}

// PlanStep is a single step in the legacy execution plan.
type PlanStep struct {
	Description string `json:"description"`
	SQL         string `json:"sql,omitempty"`
}

// NodeOutput is the result from executing a graph node.
type NodeOutput struct {
	NextNode string
	Error    error
}

// Node defines the interface for a workflow graph node.
type Node interface {
	Name() string
	Execute(ctx context.Context, state *NL2SQLState) (*NodeOutput, error)
}

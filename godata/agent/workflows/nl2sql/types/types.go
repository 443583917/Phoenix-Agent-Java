package types

import (
	"context"
	"time"
)

// NL2SQLState holds state flowing through the NL2SQL graph nodes.
type NL2SQLState struct {
	AgentID              string // originating agent
	ThreadID             string // thread / session
	Query                string // original user query
	HumanFeedback        bool   // human feedback requested
	HumanFeedbackContent string // human feedback content
	RejectedPlan         bool   // was the plan rejected
	NL2SQLOnly           bool   // NL2SQL only mode (no chat)

	Intent           string  // classified intent (e.g., "sql", "qa", "chat")
	IntentConfidence float64 // intent confidence score

	EvidenceContexts []EvidenceContext // recalled evidence/knowledge
	SchemaContext    []SchemaContext   // recalled schema info

	Plan            *ExecutionPlan           // generated plan
	SQLQuery        string                   // generated SQL
	PythonCode      string                   // generated Python (if needed)
	ExecutionResult map[string]interface{}   // execution output

	ReportContent string    // final HTML report
	Errors        []string  // accumulated errors
	CurrentNode   string    // current node name
	StartTime     time.Time // graph start time
}

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

// ExecutionPlan represents the generated plan for SQL execution.
type ExecutionPlan struct {
	Steps     []PlanStep `json:"steps"`
	Reasoning string     `json:"reasoning"`
}

// PlanStep is a single step in the execution plan.
type PlanStep struct {
	Description string `json:"description"`
	SQL         string `json:"sql,omitempty"`
}

// NodeOutput is the result from executing a graph node.
type NodeOutput struct {
	NextNode string // name of the next node (empty = end)
	Error    error  // execution error, if any
}

// Node defines the interface for a workflow graph node.
type Node interface {
	Name() string
	Execute(ctx context.Context, state *NL2SQLState) (*NodeOutput, error)
}

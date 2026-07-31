package nl2sql

import (
	"context"
	"reflect"
	"strings"

	"github.com/phoenix-agent-go/agent/knowledge"
	"github.com/phoenix-agent-go/agent/tools/datasource"
	"github.com/phoenix-agent-go/agent/workflows/nl2sql/nodes"
	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// ──────────────────────────── Node Names ────────────────────────────

const (
	NodeIntentRecognition     = "intent_recognition"
	NodeEvidenceRecall        = "evidence_recall"
	NodeQueryEnhance          = "query_enhance"
	NodeSchemaRecall          = "schema_recall"
	NodeTableRelation         = "table_relation"
	NodeFeasibilityAssessment = "feasibility_assessment"
	NodePlanner               = "planner"
	NodePlanExecutor          = "plan_executor"
	NodeSQLGenerate           = "sql_generate"
	NodeSemanticConsistency   = "semantic_consistency"
	NodeSQLExecute            = "sql_execute"
	NodePythonGenerate        = "python_generate_node"
	NodePythonExecute         = "python_execute"
	NodePythonAnalyze         = "python_analyze"
	NodeHumanFeedback         = "human_feedback"
	NodeReportGenerate        = "report_generate"
)

// ──────────────────────────── NL2SQLGraph Builder ────────────────────────────

// NL2SQLGraph builds the NL2SQL workflow graph using the tRPC-Agent-Go
// graph.StateGraph builder with all 16 nodes and complex conditional routing.
type NL2SQLGraph struct {
	llmService  *nodes.LLMService
	retriever   *knowledge.Retriever
	dbManager   *datasource.DatasourceManager
	humanReview bool
}

// NewNL2SQLGraph creates a new NL2SQLGraph builder with dependencies.
func NewNL2SQLGraph(m model.Model, retriever *knowledge.Retriever, dbManager *datasource.DatasourceManager) *NL2SQLGraph {
	return &NL2SQLGraph{
		llmService: nodes.NewLLMService(m),
		retriever:  retriever,
		dbManager:  dbManager,
	}
}

// WithHumanReview enables human-in-the-loop review.
func (g *NL2SQLGraph) WithHumanReview(enabled bool) *NL2SQLGraph {
	g.humanReview = enabled
	return g
}

// Build instantiates all 16 node stubs and wires them into a tRPC-Agent-Go
// StateGraph with complex conditional routing. Returns the compiled graph.
func (g *NL2SQLGraph) Build() (*graph.Graph, error) {
	schema := graph.NewStateSchema()
	schema.AddField("nl2sql_state", graph.StateField{
		Type:    reflectType[*nl2sqltypes.NL2SQLState](),
		Reducer: graph.DefaultReducer,
		Default: func() any { return &nl2sqltypes.NL2SQLState{} },
	})

	sg := graph.NewStateGraph(schema)

	// ── Create all 16 nodes ──
	intentNode := &nodes.IntentRecognitionNode{LLM: g.llmService}
	evidenceNode := &nodes.EvidenceRecallNode{LLM: g.llmService, Retriever: g.retriever}
	queryEnhanceNode := &nodes.QueryEnhanceNode{LLM: g.llmService}
	schemaNode := &nodes.SchemaRecallNode{Retriever: g.retriever, LLM: g.llmService}
	tableRelationNode := &nodes.TableRelationNode{LLM: g.llmService}
	feasibilityNode := &nodes.FeasibilityAssessmentNode{LLM: g.llmService}
	plannerNode := &nodes.PlannerNode{LLM: g.llmService}
	planExecNode := &nodes.PlanExecutorNode{}
	sqlGenNode := &nodes.SqlGenerateNode{LLM: g.llmService}
	semanticNode := &nodes.SemanticConsistencyNode{LLM: g.llmService}
	sqlExecNode := &nodes.SqlExecuteNode{DBManager: g.dbManager}
	pythonGenNode := &nodes.PythonGenerateNode{LLM: g.llmService}
	pythonExecNode := &nodes.PythonExecuteNode{}
	pythonAnalyzeNode := &nodes.PythonAnalyzeNode{LLM: g.llmService}
	humanFeedbackNode := &nodes.HumanFeedbackNode{}
	reportNode := &nodes.ReportGeneratorNode{LLM: g.llmService}

	// ── Register all nodes ──
	sg.AddNode(intentNode.Name(), intentNode.Execute)
	sg.AddNode(evidenceNode.Name(), evidenceNode.Execute)
	sg.AddNode(queryEnhanceNode.Name(), queryEnhanceNode.Execute)
	sg.AddNode(schemaNode.Name(), schemaNode.Execute)
	sg.AddNode(tableRelationNode.Name(), tableRelationNode.Execute)
	sg.AddNode(feasibilityNode.Name(), feasibilityNode.Execute)
	sg.AddNode(plannerNode.Name(), plannerNode.Execute)
	sg.AddNode(planExecNode.Name(), planExecNode.Execute)
	sg.AddNode(sqlGenNode.Name(), sqlGenNode.Execute)
	sg.AddNode(semanticNode.Name(), semanticNode.Execute)
	sg.AddNode(sqlExecNode.Name(), sqlExecNode.Execute)
	sg.AddNode(pythonGenNode.Name(), pythonGenNode.Execute)
	sg.AddNode(pythonExecNode.Name(), pythonExecNode.Execute)
	sg.AddNode(pythonAnalyzeNode.Name(), pythonAnalyzeNode.Execute)
	sg.AddNode(humanFeedbackNode.Name(), humanFeedbackNode.Execute)
	sg.AddNode(reportNode.Name(), reportNode.Execute)

	// ── Set entry point ──
	sg.SetEntryPoint(intentNode.Name())

	// ── Linear edges ──
	// sg.AddEdge(intentNode.Name(), evidenceNode.Name()) -- routed via conditional edge
	sg.AddEdge(evidenceNode.Name(), queryEnhanceNode.Name())
	sg.AddConditionalEdges(queryEnhanceNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state != nil && nl2state.QueryEnhanceOutput != nil && nl2state.QueryEnhanceOutput.CanonicalQuery != "" {
				return schemaNode.Name(), nil
			}
			return graph.End, nil
		},
		map[string]string{
			schemaNode.Name(): schemaNode.Name(),
			graph.End:         graph.End,
		},
	)

	sg.AddConditionalEdges(schemaNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state != nil && len(nl2state.TableDocumentsForSchema) > 0 {
				return tableRelationNode.Name(), nil
			}
			return graph.End, nil
		},
		map[string]string{
			tableRelationNode.Name(): tableRelationNode.Name(),
			graph.End:                graph.End,
		},
	)
	sg.AddEdge(plannerNode.Name(), planExecNode.Name())
	sg.AddEdge(pythonGenNode.Name(), pythonExecNode.Name())
	sg.AddEdge(pythonExecNode.Name(), pythonAnalyzeNode.Name())
	sg.AddEdge(reportNode.Name(), graph.End)

	// ── Conditional edge: table_relation -> feasibility (with self-loop retry) ──
	sg.AddConditionalEdges(tableRelationNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state != nil && nl2state.TableRelationException != "" &&
				// TODO: use cfg.MaxSQLRetryCount (default 10)
				nl2state.TableRelationRetryCount < 3 &&
				isRetryableError(nl2state.TableRelationException) {
				nl2state.TableRelationRetryCount++
				nl2state.TableRelationException = ""
				return tableRelationNode.Name(), nil
			}
			return feasibilityNode.Name(), nil
		},
		map[string]string{
			tableRelationNode.Name(): tableRelationNode.Name(),
			feasibilityNode.Name():   feasibilityNode.Name(),
		},
	)

	// ── Conditional edge 2: feasibility_assessment -> planner or end ──
	sg.AddConditionalEdges(feasibilityNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state == nil {
				return plannerNode.Name(), nil
			}
			feasibility := nl2state.FeasibilityOutput
			if feasibility == "自由闲聊" || feasibility == "需要澄清" {
				return graph.End, nil
			}
			return plannerNode.Name(), nil
		},
		map[string]string{
			plannerNode.Name(): plannerNode.Name(),
			graph.End:          graph.End,
		},
	)

	// ── Conditional edge 3: plan_executor -> next node based on plan step ──
	sg.AddConditionalEdges(planExecNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state == nil || nl2state.PlanNextNode == "" {
				return sqlGenNode.Name(), nil
			}
			return nl2state.PlanNextNode, nil
		},
		map[string]string{
			sqlGenNode.Name():    sqlGenNode.Name(),
			pythonGenNode.Name(): pythonGenNode.Name(),
			reportNode.Name():    reportNode.Name(),
			graph.End:            graph.End,
		},
	)

	// ── Conditional edge 4: sql_generate -> semantic_consistency ──
	sg.AddEdge(sqlGenNode.Name(), semanticNode.Name())

	// ── Conditional edge 5: semantic_consistency -> sql_execute or sql_generate (retry) ──
	sg.AddConditionalEdges(semanticNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state == nil {
				return sqlExecNode.Name(), nil
			}

			consistency := nl2state.SemanticConsistencyOutput
			generateCount := nl2state.SQLGenerateCount

			// If validation failed AND we haven't retried too many times, go back to regenerate
			// TODO: use cfg.MaxSQLRetryCount (default 10)
			if strings.HasPrefix(consistency, "不通过") && generateCount < 3 {
				// Store retry reason
				reason := strings.TrimPrefix(consistency, "不通过。")
				nl2state.SQLRegenReason = &nl2sqltypes.SqlRetryDto{
					Reason:       reason,
					SemanticFail: true,
				}
				return sqlGenNode.Name(), nil
			}

			return sqlExecNode.Name(), nil
		},
		map[string]string{
			sqlGenNode.Name(): sqlGenNode.Name(),
			sqlExecNode.Name(): sqlExecNode.Name(),
		},
	)

	// ── Conditional edge 6: sql_execute -> python_generate or report (if failed) ──
	sg.AddConditionalEdges(sqlExecNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state == nil {
				return pythonGenNode.Name(), nil
			}

			executeOutput := nl2state.SQLExecuteOutput
			generateCount := nl2state.SQLGenerateCount

			// If execution failed AND we haven't retried too many times, retry SQL
			// TODO: use cfg.MaxSQLRetryCount (default 10)
			if (strings.Contains(executeOutput, "失败") || executeOutput == "") && generateCount < 3 {
				nl2state.SQLRegenReason = &nl2sqltypes.SqlRetryDto{
					Reason:         "SQL执行失败: " + executeOutput,
					SQLExecuteFail: true,
				}
				return sqlGenNode.Name(), nil
			}

			// Check if we should proceed to Python analysis
			if nl2state.PlannerOutput != nil && len(nl2state.PlannerOutput.ExecutionPlan) > 1 {
				return pythonGenNode.Name(), nil
			}

			// If only NL2SQL mode or plan completed, go to report
			return reportNode.Name(), nil
		},
		map[string]string{
			sqlGenNode.Name():    sqlGenNode.Name(),
			pythonGenNode.Name(): pythonGenNode.Name(),
			reportNode.Name():    reportNode.Name(),
		},
	)

	// ── Conditional edge 7: python_analyze -> human_feedback or report ──
	sg.AddConditionalEdges(pythonAnalyzeNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state != nil && nl2state.HumanReviewEnabled {
				return humanFeedbackNode.Name(), nil
			}
			return reportNode.Name(), nil
		},
		map[string]string{
			humanFeedbackNode.Name(): humanFeedbackNode.Name(),
			reportNode.Name():        reportNode.Name(),
		},
	)

	// ── Conditional edge 8: human_feedback -> report or end ──
	sg.AddConditionalEdges(humanFeedbackNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state != nil && nl2state.HumanFeedbackData != "" {
				return reportNode.Name(), nil
			}
			// If no feedback data and human review is enabled, this means
			// the interrupt was triggered - end the graph (will resume later)
			return graph.End, nil
		},
		map[string]string{
			reportNode.Name(): reportNode.Name(),
			graph.End:         graph.End,
		},
	)

	// ── Conditional edge 9: Plann validation cycle ──
	// If planner output is invalid, route back to planner for repair
	// This is handled within the planner -> plan_executor edge
	sg.AddConditionalEdges(plannerNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state == nil {
				return planExecNode.Name(), nil
			}
			// Check repair count
			if nl2state.PlanRepairCount > 3 {
				// Too many retries, go to report with whatever we have
				return reportNode.Name(), nil
			}
			return planExecNode.Name(), nil
		},
		map[string]string{
			planExecNode.Name(): planExecNode.Name(),
			reportNode.Name():   reportNode.Name(),
		},
	)

	// ── Conditional edge 10: Intent routing ──
	// If intent is not data analysis, skip to end
	sg.AddConditionalEdges(intentNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			if nl2state != nil && nl2state.IntentOutput != nil {
				if nl2state.IntentOutput.Classification == "闲聊或无关指令" {
					return graph.End, nil
				}
			}
			return evidenceNode.Name(), nil
		},
		map[string]string{
			evidenceNode.Name(): evidenceNode.Name(),
			graph.End:           graph.End,
		},
	)

	// ── Conditional edge 11: Python execute retry ──
	// If Python execution failed, route back to Python generate for retry
	sg.AddConditionalEdges(pythonExecNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state := getState(state)
			// TODO: use cfg.PythonMaxTriesCount (default 5)
			if nl2state != nil && !nl2state.PythonIsSuccess && nl2state.PythonTriesCount < 3 {
				// Retry Python with different code
				return pythonGenNode.Name(), nil
			}
			return pythonAnalyzeNode.Name(), nil
		},
		map[string]string{
			pythonGenNode.Name():    pythonGenNode.Name(),
			pythonAnalyzeNode.Name(): pythonAnalyzeNode.Name(),
		},
	)

	return sg.Compile()
}

// ──────────────────────────── State access helpers ────────────────────────────

// getState retrieves the NL2SQLState from graph state.
func getState(state graph.State) *nl2sqltypes.NL2SQLState {
	if s, ok := state["nl2sql_state"].(*nl2sqltypes.NL2SQLState); ok && s != nil {
		return s
	}
	return nil
}

// reflectType returns the reflect.Type for a generic type T.
func reflectType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// compile-time checks
var (
	_ = NewNL2SQLGraph
)

func isRetryableError(err string) bool {
	return strings.Contains(err, "timeout") || strings.Contains(err, "连接") || strings.Contains(err, "超时") || strings.Contains(err, "retry")
}

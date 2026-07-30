// Package datasource provides function tools for executing SQL queries against
// configured datasources. It matches the Java BpmToolSearch pattern of providing
// LLM-friendly query, count, and table-listing capabilities with read-only safety.
package datasource

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/phoenix-agent-go/agent/tools"
)

const (
	maxQueryRows   = 100
	queryTimeout   = 30 * time.Second
	maxRowPreview  = 200 // characters per value in preview
)

// DatasourceManager holds named *sql.DB connections, thread-safe for concurrent use.
//
// Usage:
//
//	mgr := NewDatasourceManager()
//	mgr.Add("analytics", "postgres://user:pass@host:5432/db?sslmode=disable", "postgres")
//	mgr.Add("crm", "mysql://user:pass@tcp(host:3306)/db", "mysql")
//
// The tool functions look up datasources by the name registered here.
type DatasourceManager struct {
	mu  sync.RWMutex
	dbs map[string]*sql.DB
}

// NewDatasourceManager creates an empty DatasourceManager.
func NewDatasourceManager() *DatasourceManager {
	return &DatasourceManager{
		dbs: make(map[string]*sql.DB),
	}
}

// Add opens a new database connection and registers it under the given name.
//
// Parameters:
//   - name: logical datasource name used by the tool functions
//   - dsn:  data source name (driver-specific connection string)
//   - driver: "postgres" or "mysql"
func (m *DatasourceManager) Add(name, dsn, driver string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("datasource %q: failed to open connection: %w", name, err)
	}

	// Configure connection pool with conservative defaults.
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connectivity.
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return fmt.Errorf("datasource %q: ping failed: %w", name, err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Close existing connection if re-adding under the same name.
	if prev, ok := m.dbs[name]; ok {
		prev.Close()
	}
	m.dbs[name] = db

	zap.L().Info("datasource registered", zap.String("name", name), zap.String("driver", driver))
	return nil
}

// Remove closes and removes a datasource by name.
func (m *DatasourceManager) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if db, ok := m.dbs[name]; ok {
		db.Close()
		delete(m.dbs, name)
	}
}

// get returns the *sql.DB for the named datasource, or an error.
func (m *DatasourceManager) get(name string) (*sql.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, ok := m.dbs[name]
	if !ok {
		return nil, fmt.Errorf("datasource %q not found", name)
	}
	return db, nil
}

// ListNames returns all registered datasource names.
func (m *DatasourceManager) ListNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.dbs))
	for name := range m.dbs {
		names = append(names, name)
	}
	return names
}

// ──────────────────────────── Query Helpers ────────────────────────────

// validateReadSQL ensures the query is a read-only SELECT statement.
// Returns the trimmed SQL for execution.
func validateReadSQL(raw string) (string, error) {
	sql := strings.TrimSpace(raw)
	if sql == "" {
		return "", fmt.Errorf("SQL query must not be empty")
	}
	upper := strings.ToUpper(sql)
	if !strings.HasPrefix(upper, "SELECT") {
		return "", fmt.Errorf("only SELECT queries are allowed for read-only operations")
	}
	// Reject statements that modify data even if they start with SELECT.
	// This catches UNION-based injection attempts.
	for _, keyword := range []string{"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "CREATE", "TRUNCATE", "EXEC", "EXECUTE"} {
		if strings.Contains(upper, keyword) {
			return "", fmt.Errorf("query contains forbidden keyword: %s", keyword)
		}
	}
	return sql, nil
}

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, queryTimeout)
}

// queryRows executes a read-only SELECT query and returns columns + rows.
func (m *DatasourceManager) queryRows(ctx context.Context, datasourceName, rawSQL string) (*QueryResult, error) {
	db, err := m.get(datasourceName)
	if err != nil {
		return nil, err
	}

	query, err := validateReadSQL(rawSQL)
	if err != nil {
		return nil, err
	}

	qCtx, cancel := withTimeout(ctx)
	defer cancel()

	tx, err := db.BeginTx(qCtx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("begin read-only transaction: %w", err)
	}
	defer tx.Rollback() // always rollback; no writes

	rows, err := tx.QueryContext(qCtx, query)
	if err != nil {
		return nil, fmt.Errorf("query execution: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("get columns: %w", err)
	}

	result := &QueryResult{
		Columns: columns,
		Rows:    make([][]any, 0),
	}

	rowCount := 0
	for rows.Next() {
		if rowCount >= maxQueryRows {
			break
		}

		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		// Marshal scanned values to JSON-compatible types.
		row := make([]any, len(columns))
		for i, v := range values {
			row[i] = normalizeValue(v)
		}

		result.Rows = append(result.Rows, row)
		rowCount++
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	result.RowCount = rowCount
	result.Truncated = rowCount >= maxQueryRows

	return result, nil
}

// normalizeValue converts driver-specific types to JSON-safe equivalents.
func normalizeValue(v any) any {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case []byte:
		// Convert byte slices to strings.
		s := string(val)
		if len(s) > maxRowPreview {
			return s[:maxRowPreview] + "..."
		}
		return s
	case time.Time:
		return val.Format(time.RFC3339)
	default:
		return v
	}
}

// ──────────────────────────── Tool Input/Output Types ────────────────────────────

// QueryInput is the input for the query_datasource tool.
type QueryInput struct {
	Datasource string `json:"datasource" jsonschema:"description=Name of the datasource to query,required"`
	SQL        string `json:"sql"        jsonschema:"description=SELECT query to execute (read-only),required"`
}

// QueryResult is the structured result returned by query_datasource.
type QueryResult struct {
	Columns   []string `json:"columns"`
	Rows      [][]any  `json:"rows"`
	RowCount  int      `json:"rowCount"`
	Truncated bool     `json:"truncated,omitempty"`
}

// QueryOutput wraps QueryResult for the function tool return.
type QueryOutput struct {
	Result *QueryResult `json:"result"`
	Error  string       `json:"error,omitempty"`
}

// CountInput is the input for the count_datasource tool.
type CountInput struct {
	Datasource string `json:"datasource" jsonschema:"description=Name of the datasource to query,required"`
	SQL        string `json:"sql"        jsonschema:"description=SELECT COUNT query to execute,required"`
}

// CountOutput is the return type for count_datasource.
type CountOutput struct {
	Count int64  `json:"count"`
	Error string `json:"error,omitempty"`
}

// ListTablesInput is the input for the list_tables tool.
type ListTablesInput struct {
	Datasource string `json:"datasource" jsonschema:"description=Name of the datasource,required"`
}

// ListTablesOutput is the return type for list_tables.
type ListTablesOutput struct {
	Tables []string `json:"tables"`
	Error  string   `json:"error,omitempty"`
}

// ──────────────────────────── Tool Handler Functions ────────────────────────────

// queryHandler executes a SELECT query against a named datasource.
func (m *DatasourceManager) queryHandler(ctx context.Context, input QueryInput) (QueryOutput, error) {
	result, err := m.queryRows(ctx, input.Datasource, input.SQL)
	if err != nil {
		return QueryOutput{Error: err.Error()}, nil
	}
	return QueryOutput{Result: result}, nil
}

// countHandler executes a SELECT COUNT query and returns a single number.
func (m *DatasourceManager) countHandler(ctx context.Context, input CountInput) (CountOutput, error) {
	db, err := m.get(input.Datasource)
	if err != nil {
		return CountOutput{Error: err.Error()}, nil
	}

	query, err := validateReadSQL(input.SQL)
	if err != nil {
		return CountOutput{Error: err.Error()}, nil
	}

	qCtx, cancel := withTimeout(ctx)
	defer cancel()

	tx, err := db.BeginTx(qCtx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return CountOutput{Error: fmt.Sprintf("transaction: %v", err)}, nil
	}
	defer tx.Rollback()

	var count int64
	if err := tx.QueryRowContext(qCtx, query).Scan(&count); err != nil {
		return CountOutput{Error: fmt.Sprintf("count query: %v", err)}, nil
	}

	return CountOutput{Count: count}, nil
}

// listTablesHandler returns all tables in the public schema (Postgres) or default
// database (MySQL) for the given datasource.
func (m *DatasourceManager) listTablesHandler(ctx context.Context, input ListTablesInput) (ListTablesOutput, error) {
	db, err := m.get(input.Datasource)
	if err != nil {
		return ListTablesOutput{Error: err.Error()}, nil
	}

	qCtx, cancel := withTimeout(ctx)
	defer cancel()

	tx, err := db.BeginTx(qCtx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return ListTablesOutput{Error: fmt.Sprintf("transaction: %v", err)}, nil
	}
	defer tx.Rollback()

	// Use a driver-independent approach: try information_schema first,
	// fall back to driver-specific queries if needed.
	query := `SELECT table_name FROM information_schema.tables
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		ORDER BY table_name`

	rows, err := tx.QueryContext(qCtx, query)
	if err != nil {
		// Fallback: try SHOW TABLES for MySQL.
		rows2, err2 := tx.QueryContext(qCtx, "SHOW TABLES")
		if err2 != nil {
			return ListTablesOutput{Error: fmt.Sprintf("list tables: %v / fallback: %v", err, err2)}, nil
		}
		rows = rows2
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return ListTablesOutput{Error: fmt.Sprintf("scan: %v", err)}, nil
		}
		tables = append(tables, tableName)
	}

	if tables == nil {
		tables = []string{}
	}
	return ListTablesOutput{Tables: tables}, nil
}

// ──────────────────────────── Tool Constructors ────────────────────────────

// QueryDatasourceTool returns a function tool for executing SELECT queries.
func (m *DatasourceManager) QueryDatasourceTool() tool.Tool {
	return function.NewFunctionTool(
		m.queryHandler,
		function.WithName("query_datasource"),
		function.WithDescription(
			"Execute a read-only SELECT query against a named datasource. "+
				"Returns up to 100 rows as a JSON array with column names and row data. "+
				"Parameters: datasource (string) — name of the configured datasource; "+
				"sql (string) — the SELECT query to execute. "+
				"Use this to explore data, verify assumptions, and answer analytical questions.",
		),
	)
}

// CountDatasourceTool returns a function tool for executing COUNT queries.
func (m *DatasourceManager) CountDatasourceTool() tool.Tool {
	return function.NewFunctionTool(
		m.countHandler,
		function.WithName("count_datasource"),
		function.WithDescription(
			"Execute a SELECT COUNT query against a named datasource and return a single number. "+
				"Parameters: datasource (string) — name of the configured datasource; "+
				"sql (string) — a SELECT COUNT(...) query. "+
				"Use this before querying to understand data volume.",
		),
	)
}

// ListTablesTool returns a function tool for listing available tables.
func (m *DatasourceManager) ListTablesTool() tool.Tool {
	return function.NewFunctionTool(
		m.listTablesHandler,
		function.WithName("list_tables"),
		function.WithDescription(
			"List all tables in a named datasource. "+
				"Parameters: datasource (string) — name of the configured datasource. "+
				"Returns an array of table name strings. "+
				"Use this to discover the available data schema before writing queries.",
		),
	)
}

// Tools returns all three datasource tools as a slice.
func (m *DatasourceManager) Tools() []tool.Tool {
	return []tool.Tool{
		m.QueryDatasourceTool(),
		m.CountDatasourceTool(),
		m.ListTablesTool(),
	}
}

// ──────────────────────────── Registration Helper ────────────────────────────

// RegisterTools registers all datasource tools into the given ToolRegistry.
// This is the primary integration point for wiring tools into an agent.
//
// Usage:
//
//	mgr := datasource.NewDatasourceManager()
//	mgr.Add("analytics", "...", "postgres")
//	registry := tools.NewToolRegistry()
//	datasource.RegisterTools(registry, mgr)
func RegisterTools(registry *tools.ToolRegistry, mgr *DatasourceManager) {
	for _, t := range mgr.Tools() {
		registry.Register(t)
	}
}

// compile-time check that we satisfy the intended interface
var (
	_ json.Marshaler = (*QueryResult)(nil)
)

// MarshalJSON is a no-op; the types are serialized by the function tool framework
// automatically. Included to satisfy the json.Marshaler interface check.
func (q *QueryResult) MarshalJSON() ([]byte, error) {
	type alias QueryResult
	return json.Marshal((*alias)(q))
}

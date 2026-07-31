package config

type GraphConfig struct {
	MaxSQLRetryCount           int     `mapstructure:"max_sql_retry_count"`
	MaxSQLOptimizeCount        int     `mapstructure:"max_sql_optimize_count"`
	SQLScoreThreshold          float64 `mapstructure:"sql_score_threshold"`
	MaxTurnHistory             int     `mapstructure:"max_turn_history"`
	MaxPlanLength              int     `mapstructure:"max_plan_length"`
	MaxColumnsPerTable         int     `mapstructure:"max_columns_per_table"`
	EnableSQLResultChart       bool    `mapstructure:"enable_sql_result_chart"`
	PythonMaxTriesCount        int     `mapstructure:"python_max_tries_count"`
	TableTopkLimit             int     `mapstructure:"table_topk_limit"`
	TableSimilarityThreshold   float64 `mapstructure:"table_similarity_threshold"`
	DefaultSimilarityThreshold float64 `mapstructure:"default_similarity_threshold"`
	DefaultTopkLimit           int     `mapstructure:"default_topk_limit"`
}

func DefaultGraphConfig() GraphConfig {
	return GraphConfig{
		MaxSQLRetryCount:           10,
		MaxSQLOptimizeCount:        10,
		SQLScoreThreshold:          0.95,
		MaxTurnHistory:             5,
		MaxPlanLength:              2000,
		MaxColumnsPerTable:         150,
		EnableSQLResultChart:       true,
		PythonMaxTriesCount:        5,
		TableTopkLimit:             10,
		TableSimilarityThreshold:   0.2,
		DefaultSimilarityThreshold: 0.4,
		DefaultTopkLimit:           8,
	}
}

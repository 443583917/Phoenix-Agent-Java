package config

type VectorStoreConfig struct {
	Dimensions          int     `mapstructure:"dimensions"`
	SimilarityThreshold float64 `mapstructure:"similarity_threshold"`
	TableTopkLimit      int     `mapstructure:"table_topk_limit"`
	DefaultTopkLimit    int     `mapstructure:"default_topk_limit"`
	BatchDelTopkLimit   int     `mapstructure:"batch_del_topk_limit"`
	EnableHybridSearch  bool    `mapstructure:"enable_hybrid_search"`
}

func DefaultVectorStoreConfig() VectorStoreConfig {
	return VectorStoreConfig{
		Dimensions:          512,
		SimilarityThreshold: 0.4,
		TableTopkLimit:      10,
		DefaultTopkLimit:    8,
		BatchDelTopkLimit:   5000,
		EnableHybridSearch:  false,
	}
}

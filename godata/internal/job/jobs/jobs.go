package jobs

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AgentStatisticsJob struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAgentStatisticsJob(db *gorm.DB) *AgentStatisticsJob {
	return &AgentStatisticsJob{
		db:     db,
		logger: zap.L().Named("job.agent_statistics"),
	}
}

func (j *AgentStatisticsJob) Name() string {
	return "agent_statistics"
}

func (j *AgentStatisticsJob) Run(ctx context.Context) error {
	j.logger.Info("calculating agent statistics")

	var count int64
	if err := j.db.WithContext(ctx).Table("tbl_data_agent").Where("del_flag = 0").Count(&count).Error; err != nil {
		return err
	}

	j.logger.Info("agent statistics calculated", zap.Int64("totalAgents", count))
	return nil
}

type KnowledgeEmbeddingRetryJob struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewKnowledgeEmbeddingRetryJob(db *gorm.DB) *KnowledgeEmbeddingRetryJob {
	return &KnowledgeEmbeddingRetryJob{
		db:     db,
		logger: zap.L().Named("job.knowledge_embedding_retry"),
	}
}

func (j *KnowledgeEmbeddingRetryJob) Name() string {
	return "knowledge_embedding_retry"
}

func (j *KnowledgeEmbeddingRetryJob) Run(ctx context.Context) error {
	j.logger.Info("checking for failed embeddings to retry")

	result := j.db.WithContext(ctx).
		Table("tbl_data_agent_knowledge").
		Where("embedding_status = ? AND del_flag = 0", "failed").
		Update("embedding_status", "pending")

	if result.Error != nil {
		return result.Error
	}

	j.logger.Info("embedding retry check completed", zap.Int64("retried", result.RowsAffected))
	return nil
}

type SessionCleanupJob struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewSessionCleanupJob(db *gorm.DB) *SessionCleanupJob {
	return &SessionCleanupJob{
		db:     db,
		logger: zap.L().Named("job.session_cleanup"),
	}
}

func (j *SessionCleanupJob) Name() string {
	return "session_cleanup"
}

func (j *SessionCleanupJob) Run(ctx context.Context) error {
	j.logger.Info("cleaning up old sessions")

	cutoff := time.Now().AddDate(0, -3, 0)
	result := j.db.WithContext(ctx).
		Table("tbl_data_chat_session").
		Where("status = ? AND update_time < ? AND del_flag = 0", "active", cutoff).
		Update("status", "closed")

	if result.Error != nil {
		return result.Error
	}

	j.logger.Info("session cleanup completed", zap.Int64("closed", result.RowsAffected))
	return nil
}

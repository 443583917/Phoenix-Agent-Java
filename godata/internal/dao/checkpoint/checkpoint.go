package checkpoint

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Checkpoint struct {
	ID        string    `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	ThreadID  string    `gorm:"column:thread_id;type:varchar(64);index" json:"threadId"`
	StateJSON string    `gorm:"column:state_json;type:text" json:"stateJson"`
	NodeName  string    `gorm:"column:node_name;type:varchar(128)" json:"nodeName"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (Checkpoint) TableName() string { return "tbl_graph_checkpoint" }

type CheckpointStore struct {
	db *gorm.DB
}

func NewCheckpointStore(db *gorm.DB) *CheckpointStore {
	return &CheckpointStore{db: db}
}

func (s *CheckpointStore) Save(ctx context.Context, threadID, node string, state map[string]interface{}) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	cp := &Checkpoint{
		ID:        threadID + "_" + node,
		ThreadID:  threadID,
		StateJSON: string(data),
		NodeName:  node,
	}
	return s.db.WithContext(ctx).Save(cp).Error
}

func (s *CheckpointStore) Load(ctx context.Context, threadID string) (map[string]interface{}, error) {
	var cp Checkpoint
	err := s.db.WithContext(ctx).Where("thread_id = ?", threadID).Order("created_at DESC").First(&cp).Error
	if err != nil {
		return nil, err
	}
	var state map[string]interface{}
	if err := json.Unmarshal([]byte(cp.StateJSON), &state); err != nil {
		return nil, err
	}
	return state, nil
}

func (s *CheckpointStore) Delete(ctx context.Context, threadID string) error {
	return s.db.WithContext(ctx).Where("thread_id = ?", threadID).Delete(&Checkpoint{}).Error
}

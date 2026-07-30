package db

import (
	"context"
	"errors"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"gorm.io/gorm"
)

// ──────────────────────────── GroupInfo ────────────────────────────

type groupInfoRepo struct{ db *gorm.DB }

func NewGroupInfoRepository(db *gorm.DB) repository.GroupInfoRepository {
	return &groupInfoRepo{db}
}

func (r *groupInfoRepo) FindByID(ctx context.Context, id string) (*model.GroupInfo, error) {
	var entity model.GroupInfo
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *groupInfoRepo) Page(ctx context.Context, page, size int, query *model.GroupInfo) ([]*model.GroupInfo, int64, error) {
	var list []*model.GroupInfo
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.GroupInfo{}).Where("del_flag = 0")
	if query != nil {
		if query.Name != "" {
			dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.SN != "" {
			dbQuery = dbQuery.Where("sn LIKE ?", "%"+query.SN+"%")
		}
		if query.Status != 0 {
			dbQuery = dbQuery.Where("status = ?", query.Status)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *groupInfoRepo) List(ctx context.Context) ([]*model.GroupInfo, error) {
	var list []*model.GroupInfo
	err := r.db.WithContext(ctx).Where("del_flag = 0").Find(&list).Error
	return list, err
}

func (r *groupInfoRepo) Create(ctx context.Context, entity *model.GroupInfo) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *groupInfoRepo) Update(ctx context.Context, entity *model.GroupInfo) error {
	return r.db.WithContext(ctx).Model(&model.GroupInfo{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *groupInfoRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.GroupInfo{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

func (r *groupInfoRepo) ToggleStatus(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.GroupInfo{}).Where("id = ? AND del_flag = 0", id).
		Update("status", gorm.Expr("CASE WHEN status = 0 THEN 1 ELSE 0 END")).Error
}

func (r *groupInfoRepo) GetEnabled(ctx context.Context) ([]*model.GroupInfo, error) {
	var list []*model.GroupInfo
	err := r.db.WithContext(ctx).Where("del_flag = 0 AND status = 0").Find(&list).Error
	return list, err
}

// ──────────────────────────── GroupAgentInfo ────────────────────────────

type groupAgentInfoRepo struct{ db *gorm.DB }

func NewGroupAgentInfoRepository(db *gorm.DB) repository.GroupAgentInfoRepository {
	return &groupAgentInfoRepo{db}
}

func (r *groupAgentInfoRepo) FindByID(ctx context.Context, id string) (*model.GroupAgentInfo, error) {
	var entity model.GroupAgentInfo
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *groupAgentInfoRepo) FindByGroupID(ctx context.Context, groupID string) ([]*model.GroupAgentInfo, error) {
	var list []*model.GroupAgentInfo
	err := r.db.WithContext(ctx).Where("group_id = ? AND del_flag = 0", groupID).Find(&list).Error
	return list, err
}

func (r *groupAgentInfoRepo) FindByAgentID(ctx context.Context, agentID int64) ([]*model.GroupAgentInfo, error) {
	var list []*model.GroupAgentInfo
	err := r.db.WithContext(ctx).Where("agent_id = ? AND del_flag = 0", agentID).Find(&list).Error
	return list, err
}

func (r *groupAgentInfoRepo) GetByGroupIdAndAgentId(ctx context.Context, groupID string, agentID int64) (*model.GroupAgentInfo, error) {
	var entity model.GroupAgentInfo
	err := r.db.WithContext(ctx).Where("group_id = ? AND agent_id = ? AND del_flag = 0", groupID, agentID).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *groupAgentInfoRepo) List(ctx context.Context) ([]*model.GroupAgentInfo, error) {
	var list []*model.GroupAgentInfo
	err := r.db.WithContext(ctx).Where("del_flag = 0").Find(&list).Error
	return list, err
}

func (r *groupAgentInfoRepo) Create(ctx context.Context, entity *model.GroupAgentInfo) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *groupAgentInfoRepo) Update(ctx context.Context, entity *model.GroupAgentInfo) error {
	return r.db.WithContext(ctx).Model(&model.GroupAgentInfo{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *groupAgentInfoRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.GroupAgentInfo{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── AccountInfo ────────────────────────────

type accountInfoRepo struct{ db *gorm.DB }

func NewAccountInfoRepository(db *gorm.DB) repository.AccountInfoRepository {
	return &accountInfoRepo{db}
}

func (r *accountInfoRepo) FindByID(ctx context.Context, id string) (*model.AccountInfo, error) {
	var entity model.AccountInfo
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *accountInfoRepo) FindByUsername(ctx context.Context, username string) (*model.AccountInfo, error) {
	var entity model.AccountInfo
	err := r.db.WithContext(ctx).Where("username = ? AND del_flag = 0", username).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *accountInfoRepo) FindByThirdPartyID(ctx context.Context, thirdPartyID string) (*model.AccountInfo, error) {
	var entity model.AccountInfo
	err := r.db.WithContext(ctx).Where("third_party_id = ? AND del_flag = 0", thirdPartyID).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *accountInfoRepo) FindByCode(ctx context.Context, code string) (*model.AccountInfo, error) {
	var entity model.AccountInfo
	err := r.db.WithContext(ctx).Where("code = ? AND del_flag = 0", code).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *accountInfoRepo) FindByStatus(ctx context.Context, status string) ([]*model.AccountInfo, error) {
	var list []*model.AccountInfo
	err := r.db.WithContext(ctx).Where("status = ? AND del_flag = 0", status).Find(&list).Error
	return list, err
}

func (r *accountInfoRepo) Page(ctx context.Context, page, size int, query *model.AccountInfo) ([]*model.AccountInfo, int64, error) {
	var list []*model.AccountInfo
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.AccountInfo{}).Where("del_flag = 0")
	if query != nil {
		if query.Username != "" {
			dbQuery = dbQuery.Where("username LIKE ?", "%"+query.Username+"%")
		}
		if query.RealName != "" {
			dbQuery = dbQuery.Where("real_name LIKE ?", "%"+query.RealName+"%")
		}
		if query.Code != "" {
			dbQuery = dbQuery.Where("code LIKE ?", "%"+query.Code+"%")
		}
		if query.Status != "" {
			dbQuery = dbQuery.Where("status = ?", query.Status)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *accountInfoRepo) GetUnGroupPage(ctx context.Context, page, size int, groupID string, query *model.AccountInfo) ([]*model.AccountInfo, int64, error) {
	var list []*model.AccountInfo
	var total int64

	// 查询未分配到指定分组的账号
	// SELECT * FROM tbl_platform_account_info WHERE del_flag = 0 AND id NOT IN (
	//   SELECT account_id FROM tbl_platform_account_group_info WHERE group_id = ? AND del_flag = 0
	// )
	subQuery := r.db.WithContext(ctx).Model(&model.AccountGroupInfo{}).
		Select("account_id").
		Where("group_id = ? AND del_flag = 0", groupID)

	dbQuery := r.db.WithContext(ctx).Model(&model.AccountInfo{}).
		Where("del_flag = 0").
		Where("id NOT IN (?)", subQuery)

	if query != nil {
		if query.Username != "" {
			dbQuery = dbQuery.Where("username LIKE ?", "%"+query.Username+"%")
		}
		if query.RealName != "" {
			dbQuery = dbQuery.Where("real_name LIKE ?", "%"+query.RealName+"%")
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *accountInfoRepo) GetMyAgents(ctx context.Context, accountID string) ([]*model.GroupAgentInfo, error) {
	// 通过 account_group_info 找到该账号所属的 group_id
	// 再通过 group_agent_info 找到这些 group 关联的 agent
	var groupIDs []string
	if err := r.db.WithContext(ctx).Model(&model.AccountGroupInfo{}).
		Select("group_id").
		Where("account_id = ? AND del_flag = 0", accountID).
		Find(&groupIDs).Error; err != nil {
		return nil, err
	}

	if len(groupIDs) == 0 {
		return nil, nil
	}

	var list []*model.GroupAgentInfo
	err := r.db.WithContext(ctx).Where("group_id IN ? AND del_flag = 0", groupIDs).Find(&list).Error
	return list, err
}

func (r *accountInfoRepo) List(ctx context.Context) ([]*model.AccountInfo, error) {
	var list []*model.AccountInfo
	err := r.db.WithContext(ctx).Where("del_flag = 0").Find(&list).Error
	return list, err
}

func (r *accountInfoRepo) Create(ctx context.Context, entity *model.AccountInfo) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *accountInfoRepo) Update(ctx context.Context, entity *model.AccountInfo) error {
	return r.db.WithContext(ctx).Model(&model.AccountInfo{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *accountInfoRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.AccountInfo{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

func (r *accountInfoRepo) BatchUpdateStatus(ctx context.Context, ids []string, status string) error {
	return r.db.WithContext(ctx).Model(&model.AccountInfo{}).
		Where("id IN ? AND del_flag = 0", ids).
		Update("status", status).Error
}

// ──────────────────────────── AccountGroupInfo ────────────────────────────

type accountGroupInfoRepo struct{ db *gorm.DB }

func NewAccountGroupInfoRepository(db *gorm.DB) repository.AccountGroupInfoRepository {
	return &accountGroupInfoRepo{db}
}

func (r *accountGroupInfoRepo) FindByID(ctx context.Context, id string) (*model.AccountGroupInfo, error) {
	var entity model.AccountGroupInfo
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *accountGroupInfoRepo) FindByGroupID(ctx context.Context, groupID string) ([]*model.AccountGroupInfo, error) {
	var list []*model.AccountGroupInfo
	err := r.db.WithContext(ctx).Where("group_id = ? AND del_flag = 0", groupID).Find(&list).Error
	return list, err
}

func (r *accountGroupInfoRepo) FindByAccountID(ctx context.Context, accountID string) ([]*model.AccountGroupInfo, error) {
	var list []*model.AccountGroupInfo
	err := r.db.WithContext(ctx).Where("account_id = ? AND del_flag = 0", accountID).Find(&list).Error
	return list, err
}

func (r *accountGroupInfoRepo) Page(ctx context.Context, page, size int, query *model.AccountGroupInfo) ([]*model.AccountGroupInfo, int64, error) {
	var list []*model.AccountGroupInfo
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.AccountGroupInfo{}).Where("del_flag = 0")
	if query != nil {
		if query.GroupID != "" {
			dbQuery = dbQuery.Where("group_id = ?", query.GroupID)
		}
		if query.AccountID != "" {
			dbQuery = dbQuery.Where("account_id = ?", query.AccountID)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *accountGroupInfoRepo) Create(ctx context.Context, entity *model.AccountGroupInfo) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *accountGroupInfoRepo) Update(ctx context.Context, entity *model.AccountGroupInfo) error {
	return r.db.WithContext(ctx).Model(&model.AccountGroupInfo{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *accountGroupInfoRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.AccountGroupInfo{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── AccountTenantInfo ────────────────────────────

type accountTenantInfoRepo struct{ db *gorm.DB }

func NewAccountTenantInfoRepository(db *gorm.DB) repository.AccountTenantInfoRepository {
	return &accountTenantInfoRepo{db}
}

func (r *accountTenantInfoRepo) FindByID(ctx context.Context, id string) (*model.AccountTenantInfo, error) {
	var entity model.AccountTenantInfo
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *accountTenantInfoRepo) FindByAccountID(ctx context.Context, accountID string) ([]*model.AccountTenantInfo, error) {
	var list []*model.AccountTenantInfo
	err := r.db.WithContext(ctx).Where("account_id = ? AND del_flag = 0", accountID).Find(&list).Error
	return list, err
}

func (r *accountTenantInfoRepo) FindByTenantID(ctx context.Context, tenantID string) ([]*model.AccountTenantInfo, error) {
	var list []*model.AccountTenantInfo
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND del_flag = 0", tenantID).Find(&list).Error
	return list, err
}

func (r *accountTenantInfoRepo) Page(ctx context.Context, page, size int, query *model.AccountTenantInfo) ([]*model.AccountTenantInfo, int64, error) {
	var list []*model.AccountTenantInfo
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.AccountTenantInfo{}).Where("del_flag = 0")
	if query != nil {
		if query.AccountID != "" {
			dbQuery = dbQuery.Where("account_id = ?", query.AccountID)
		}
		if query.TenantID != "" {
			dbQuery = dbQuery.Where("tenant_id = ?", query.TenantID)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *accountTenantInfoRepo) Create(ctx context.Context, entity *model.AccountTenantInfo) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *accountTenantInfoRepo) Update(ctx context.Context, entity *model.AccountTenantInfo) error {
	return r.db.WithContext(ctx).Model(&model.AccountTenantInfo{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *accountTenantInfoRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.AccountTenantInfo{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── TenantInfo ────────────────────────────

type tenantInfoRepo struct{ db *gorm.DB }

func NewTenantInfoRepository(db *gorm.DB) repository.TenantInfoRepository {
	return &tenantInfoRepo{db}
}

func (r *tenantInfoRepo) FindByID(ctx context.Context, id string) (*model.TenantInfo, error) {
	var entity model.TenantInfo
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *tenantInfoRepo) FindBySN(ctx context.Context, sn string) (*model.TenantInfo, error) {
	var entity model.TenantInfo
	err := r.db.WithContext(ctx).Where("sn = ? AND del_flag = 0", sn).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *tenantInfoRepo) Page(ctx context.Context, page, size int, query *model.TenantInfo) ([]*model.TenantInfo, int64, error) {
	var list []*model.TenantInfo
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.TenantInfo{}).Where("del_flag = 0")
	if query != nil {
		if query.Name != "" {
			dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.SN != "" {
			dbQuery = dbQuery.Where("sn LIKE ?", "%"+query.SN+"%")
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *tenantInfoRepo) Create(ctx context.Context, entity *model.TenantInfo) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *tenantInfoRepo) Update(ctx context.Context, entity *model.TenantInfo) error {
	return r.db.WithContext(ctx).Model(&model.TenantInfo{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *tenantInfoRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.TenantInfo{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── PlatformInfo ────────────────────────────

type platformInfoRepo struct{ db *gorm.DB }

func NewPlatformInfoRepository(db *gorm.DB) repository.PlatformInfoRepository {
	return &platformInfoRepo{db}
}

func (r *platformInfoRepo) FindByID(ctx context.Context, id string) (*model.PlatformInfo, error) {
	var entity model.PlatformInfo
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *platformInfoRepo) FindByType(ctx context.Context, typ string) ([]*model.PlatformInfo, error) {
	var list []*model.PlatformInfo
	err := r.db.WithContext(ctx).Where("type = ? AND del_flag = 0", typ).Find(&list).Error
	return list, err
}

func (r *platformInfoRepo) FindEnabledByType(ctx context.Context, typ string) ([]*model.PlatformInfo, error) {
	var list []*model.PlatformInfo
	err := r.db.WithContext(ctx).Where("type = ? AND status = '0' AND del_flag = 0", typ).Find(&list).Error
	return list, err
}

func (r *platformInfoRepo) Page(ctx context.Context, page, size int, query *model.PlatformInfo) ([]*model.PlatformInfo, int64, error) {
	var list []*model.PlatformInfo
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.PlatformInfo{}).Where("del_flag = 0")
	if query != nil {
		if query.Type != "" {
			dbQuery = dbQuery.Where("type = ?", query.Type)
		}
		if query.Name != "" {
			dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.Status != "" {
			dbQuery = dbQuery.Where("status = ?", query.Status)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *platformInfoRepo) Create(ctx context.Context, entity *model.PlatformInfo) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *platformInfoRepo) Update(ctx context.Context, entity *model.PlatformInfo) error {
	return r.db.WithContext(ctx).Model(&model.PlatformInfo{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *platformInfoRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.PlatformInfo{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

func (r *platformInfoRepo) ToggleStatus(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.PlatformInfo{}).Where("id = ? AND del_flag = 0", id).
		Update("status", gorm.Expr("CASE WHEN status = '0' THEN '1' ELSE '0' END")).Error
}

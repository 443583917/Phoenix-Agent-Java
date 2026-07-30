package db

import (
	"fmt"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ColumnInfo struct {
	ColumnName string `json:"columnName"`
	DataType   string `json:"dataType"`
	TableName  string `json:"tableName"`
	Comment    string `json:"comment,omitempty"`
}

type TableInfo struct {
	TableName string `json:"tableName"`
	Comment   string `json:"comment,omitempty"`
}

type datasourceAccessor struct{}

var _ repository.DatasourceAccessor = (*datasourceAccessor)(nil)

func NewDatasourceAccessor() repository.DatasourceAccessor {
	return &datasourceAccessor{}
}

func (a *datasourceAccessor) TestConnection(ds *model.Datasource) error {
	return TestDatasourceConnection(ds)
}

func (a *datasourceAccessor) GetTables(ds *model.Datasource) (interface{}, error) {
	return GetDatasourceTables(ds)
}

func (a *datasourceAccessor) GetColumns(ds *model.Datasource, tableName string) (interface{}, error) {
	return GetDatasourceColumns(ds, tableName)
}

func TestDatasourceConnection(ds *model.Datasource) error {
	db, err := openDatasourceConn(ds)
	if err != nil {
		return err
	}
	defer closeGormDB(db)
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func GetDatasourceTables(ds *model.Datasource) ([]TableInfo, error) {
	db, err := openDatasourceConn(ds)
	if err != nil {
		return nil, err
	}
	defer closeGormDB(db)

	var tables []TableInfo
	err = db.Raw(`SELECT table_name, '' as comment FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' ORDER BY table_name`).Scan(&tables).Error
	return tables, err
}

func GetDatasourceColumns(ds *model.Datasource, tableName string) ([]ColumnInfo, error) {
	db, err := openDatasourceConn(ds)
	if err != nil {
		return nil, err
	}
	defer closeGormDB(db)

	var columns []ColumnInfo
	err = db.Raw(`SELECT column_name AS columnName, data_type AS dataType, ? AS tableName FROM information_schema.columns WHERE table_schema = 'public' AND table_name = ? ORDER BY ordinal_position`, tableName, tableName).Scan(&columns).Error
	return columns, err
}

func openDatasourceConn(ds *model.Datasource) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch ds.Type {
	case "PostgreSQL":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			ds.Host, ds.Port, ds.Username, ds.Password, ds.DatabaseName)
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported datasource type for direct connection: %s", ds.Type)
	}
	return gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
}

func closeGormDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}

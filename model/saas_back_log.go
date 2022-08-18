package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type GvaModel struct {
	ID        uint           `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

// BackLog 结构体
type BackLog struct {
	GvaModel
	FinishedAt    time.Time      // 完成时间
	DomainId      int            `json:"domainId" form:"domainId" gorm:"column:domain_id;index;comment:;"`
	BackupType    string         `json:"backupType" form:"backupType" gorm:"column:backup_type;type:enum('mysqldump','xtrafull','xtraincr','mydumper','redis','tidb');comment:备份类型;"`
	DataSize      int            `json:"dataSize" form:"dataSize" gorm:"column:data_size;comment:;"`
	Status        string         `json:"status" form:"status" gorm:"column:status;type:enum('backup','success','failed');comment:备份类型;"`
	BackUpFeature *BackUpFeature `json:"backupfeature" gorm:"TYPE:json"`
}
type BackUpFeature struct {
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
}

func (c BackUpFeature) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *BackUpFeature) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

// TableName BackLog 表名
func (BackLog) TableName() string {
	return "saas_back_log"
}

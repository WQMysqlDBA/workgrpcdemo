package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateBackLog(db *gorm.DB, backlog *BackLog) error {
	return db.Debug().Create(&backlog).Error
}

// func UpdateBackLogByUuid update the backlog table by uuid
func UpdateBackLogByUuid(db *gorm.DB, uuid uuid.UUID, setvars map[string]interface{}) {
	db.Debug().Model(&BackLog{}).Where("back_up_uuid = ?", uuid).Updates(setvars)
}

// func UpdateBackLogJsonByUuid update the backlog use native sql
func UpdateBackLogJsonByUuid(db *gorm.DB, uuid uuid.UUID, status, fd string, size int64) {
	db.Debug().Model(&BackLog{}).Exec("UPDATE saas_back_log SET data_size = ? , status = ? ,back_up_feature = JSON_SET(back_up_feature,\"$.filename\", ? ) WHERE back_up_uuid = ?", size, status, fd, uuid)
}

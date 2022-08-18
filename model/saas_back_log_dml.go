package model

import "gorm.io/gorm"

func CreateBackLog(db *gorm.DB, backlog *BackLog) error {
	return db.Debug().Create(&backlog).Error
}

func UpdateBackLogById(db *gorm.DB, id int, setvars map[string]interface{}) {

}

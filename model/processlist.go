package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type InformationSchemaProcesslist struct {
	Id      int    `json:"id" form:"id" ,gorm:"column:ID;TYPE:bigint"`
	User    string `json:"user" ,from:"user" ,gorm:"column:USER;TYPE:varchar(32)"`
	Host    string `json:"host" ,form:"host" ,gorm:"column:HOST;TYPE:varchar(261)"`
	Db      string `json:"db" ,form:"db" ,gorm:"column:DB;TYPE:varchar(64)"`
	Command string `json:"command" ,form:"command" ,gorm:"column:COMMAND;TYPE:varchar(16)"`
	Time    int    `json:"time" ,form:"time" ,gorm:"column:TIME;TYPE:int"`
	State   string `json:"state" ,form:"state" ,gorm:"column:STATE;TYPE:varchar(64)"`
	Info    string `json:"info" ,form:"info" ,gorm:"column:INFO;TYPE:varchar(65535)"`
}

func (InformationSchemaProcesslist) TableName() string {
	return "PROCESSLIST"
}

func (r *InformationSchemaProcesslist) GetAllProcesslist(db *gorm.DB) []InformationSchemaProcesslist {
	// TODO 这里没有获取到结果
	resArr := make([]InformationSchemaProcesslist, 0, 0)
	db.Debug().Table("PROCESSLIST").Find(&r).Scan(&resArr)

	for _, v := range resArr {
		if data, err := json.MarshalIndent(v, "", "\t"); err == nil {
			fmt.Println(string(data))
		}
	}
	return resArr
}

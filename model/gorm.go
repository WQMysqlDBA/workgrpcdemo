package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type connInfo struct {
	User   string
	Passwd string
	Host   string
	Db     string
	Port   int
}

func newConnInfo(user, passwd, host, db string, port int) *connInfo {
	return &connInfo{
		User:   user,
		Passwd: passwd,
		Host:   host,
		Port:   port,
		Db:     db,
	}
}

func (c *connInfo) newConn() (db *gorm.DB) {
	// 用户名:密码@tcp(ip:port)/数据库?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=10s&writeTimeout=10s", c.User, c.Passwd, c.Host, c.Port, c.Db)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	// TODO 错误处理 不panic
	if err != nil {
		log.Panic("Conn mysql :%v failed,err:%v", dsn, err)
	}
	return db
}

type NewConn interface {
	newConn()
}

func GormMysql(user, passwd, host, db string, port int) (dbconn *gorm.DB) {
	c := newConnInfo(user, passwd, host, db, port)
	return c.newConn()
}

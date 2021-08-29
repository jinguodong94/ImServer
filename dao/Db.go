package dao

import (
	"fmt"
	"gindemo/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB

func InitMysql() {

	log.Println("初始数据库连接")

	config := conf.Configs.MysqlConfig

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User,
		config.Pwd, config.Address, config.Db)

	var err error

	defer func() {
		if err != nil {
			panic(fmt.Sprintf("init db error url -> %s", dsn))
		}
	}()

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	//DB, err := db.DB()
}

func CloseMysql() {
	DB, _ := Db.DB()
	DB.Close()
}

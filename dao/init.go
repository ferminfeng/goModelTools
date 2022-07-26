package dao

import (
	"log"
	"modeltools/config"
	"os"
)

func Init(dbConf *config.Database) {
	// 初始化Mysql连接池
	mysql := GetMysqlInstance().InitMysqlPool(dbConf)
	if !mysql {
		log.Println("init database pool failure...")
		os.Exit(1)
	}
}

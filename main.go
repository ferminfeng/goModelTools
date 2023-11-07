package main

import (
	"flag"
	"fmt"
	"modeltools/config"
	"modeltools/dao"
	"modeltools/generate"
)

var (
	configFile, modelPath, modelReplace, tableName, databases string
)

func main() {

	flag.StringVar(&configFile, "c", "./config.toml", "model path")

	// model保存路径
	flag.StringVar(&modelPath, "model_path", "./models/", "goModelTools server config")

	// 是否覆盖已存在model
	flag.StringVar(&modelReplace, "model_replace", "true", "model replace")

	// 待生成model的表
	flag.StringVar(&tableName, "table_name", "", "table_name")

	// 用哪个databases
	flag.StringVar(&databases, "databases", "db", "table_name")

	flag.Parse()

	// 初始化配置文件
	cfg := config.Init(configFile)

	cfg.ModelPath = modelPath
	cfg.ModelReplace = modelReplace

	dbName := ""

	// 初始化数据库
	if databases == "db" {
		dao.Init(cfg.DB)
		dbName = cfg.DB.WriteDB.DB
	} else if databases == "db_game" {
		dao.Init(cfg.DBGame)
		dbName = cfg.DBGame.WriteDB.DB
	} else {
		fmt.Println("暂不支持该数据库")
		return
	}

	if tableName == "" {
		generate.Generate(cfg, dbName) // 生成所有表信息
	} else {
		generate.Generate(cfg, dbName, tableName) // 接收传参
		// generate.Generate(cfg, "table1", "table2") // 生成指定表信息，可变参数可传入多个表名
	}
}

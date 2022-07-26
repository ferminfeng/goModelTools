package main

import (
	"flag"
	"modeltools/config"
	"modeltools/dao"
	"modeltools/generate"
)

var (
	configFile, modelPath, modelReplace string
)

func main() {

	flag.StringVar(&configFile, "c", "./config.toml", "model path")

	// model保存路径
	flag.StringVar(&modelPath, "model_path", "./models/base/", "goModelTools server config")

	// 是否覆盖已存在model
	flag.StringVar(&modelReplace, "model_replace", "true", "model replace")
	flag.Parse()

	// 初始化配置文件
	cfg := config.Init(configFile)

	cfg.ModelPath = modelPath
	cfg.ModelReplace = modelReplace

	// 初始化数据库
	dao.Init(cfg.DB)

	// generate.Generate() //生成所有表信息
	generate.Generate(cfg, "introducer", "against_task") // 生成指定表信息，可变参数可传入多个表名
}

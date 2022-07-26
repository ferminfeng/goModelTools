package conf

// ModelPath model保存路径
const ModelPath = "./models/"

// ModelReplace 是否覆盖已存在model
const ModelReplace = true

// DriverName 数据库驱动
const DriverName = "mysql"

type DbConf struct {
	Host   string
	Port   string
	User   string
	Pwd    string
	DbName string
}

// MasterDbConfig 数据库链接配置
var MasterDbConfig DbConf = DbConf{
	Host:   "127.0.0.1",
	Port:   "3306",
	User:   "root",
	Pwd:    "123456",
	DbName: "test",
}

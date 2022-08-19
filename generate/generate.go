package generate

import (
	"fmt"
	"io"
	"modeltools/config"
	"modeltools/dao"
	"modeltools/helper"
	"os"
	"strings"
)

type Table struct {
	Name    string `gorm:"column:Name"`
	Comment string `gorm:"column:Comment"`
}

type Field struct {
	Field      string `gorm:"column:Field"`
	Type       string `gorm:"column:Type"`
	Null       string `gorm:"column:Null"`
	Key        string `gorm:"column:Key"`
	Default    string `gorm:"column:Default"`
	Extra      string `gorm:"column:Extra"`
	Privileges string `gorm:"column:Privileges"`
	Comment    string `gorm:"column:Comment"`
}

// Generate 生成表
func Generate(cfg *config.Config, dbName string, tableNames ...string) {

	tableNamesStr := ""
	for _, name := range tableNames {
		if tableNamesStr != "" {
			tableNamesStr += ","
		}
		tableNamesStr += "'" + name + "'"
	}

	tables := getTables(dbName, tableNamesStr) // 生成所有表信息
	for _, table := range tables {
		fields := getFields(table.Name)
		generateModel(cfg, table, fields)
	}
}

// 获取表信息
func getTables(dbName, tableNames string) []Table {

	db := dao.GetMysqlDb()
	var tables []Table
	if tableNames == "" {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + dbName + "';").Find(&tables)
	} else {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN (" + tableNames + ") AND table_schema='" + dbName + "';").Find(&tables)
	}
	return tables
}

// 获取所有字段信息
func getFields(tableName string) []Field {
	db := dao.GetMysqlDb()
	var fields []Field
	db.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}

// 生成Model
func generateModel(cfg *config.Config, table Table, fields []Field) {
	// 获取包名
	// 指定分隔符
	countSplit := strings.Split(cfg.ModelPath, "/")
	lenth := len(countSplit)
	packageName := countSplit[lenth-2]

	// packageContent := "package models\n\n"
	packageContent := "package " + packageName + "\n\n"

	importContent := "import \"time\"\n\n"
	isUseImport := false

	tableName := helper.CamelCase(table.Name)

	tableContent := ""
	tableContent += "func (dao " + tableName + ") TableName() string {" + "\n"
	tableContent += "	return \"" + table.Name + "\"" + "\n"
	tableContent += "}" + "\n" + "\n"

	// 表注释
	if len(table.Comment) > 0 {
		tableContent += "// " + tableName + " " + table.Comment + "\n"
	}
	tableContent += "type " + tableName + " struct {\n"
	// 生成字段
	for _, field := range fields {
		// 字段名称
		fieldName := helper.CamelCase(field.Field)

		// 获取字段类型
		fieldType := getFiledType(field)

		// 获取字段json描述
		// fieldJson := getFieldJson(field)
		fieldJson := getFieldJson2(field)

		// 如果存在time类型，则引入time包
		if fieldType == "time.Time" {
			isUseImport = true
		}

		// 字段注释
		// fieldComment := getFieldComment(field)
		fieldComment := ""

		tableContent += "	" + fieldName + " " + fieldType + " `" + fieldJson + "` " + fieldComment + "\n"
	}

	tableContent += "}"

	filename := cfg.ModelPath + helper.CamelCase(table.Name) + ".go"
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		if cfg.ModelReplace != "true" {
			fmt.Println(helper.CamelCase(table.Name) + " 已存在，需删除才能重新生成...")
			return
		}
		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666) // 打开文件
		if err != nil {
			panic(err)
		}
	} else {
		f, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()

	content := packageContent
	if isUseImport {
		content += importContent + tableContent
	} else {
		content += tableContent
	}
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(helper.CamelCase(table.Name) + " 已生成...")
	}
}

// 获取字段类型
func getFiledType(field Field) string {
	typeArr := strings.Split(field.Type, "(")

	switch typeArr[0] {
	case "int":
		return "int"
	case "integer":
		return "int"
	case "mediumint":
		return "int"
	case "bit":
		return "int"
	case "year":
		return "int"
	case "smallint":
		return "int"
	case "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "decimal":
		return "float32"
	case "double":
		return "float32"
	case "float":
		return "float32"
	case "real":
		return "float32"
	case "numeric":
		return "float32"
	case "timestamp":
		return "time.Time"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	default:
		return "string"
	}
}

// 获取字段json描述
func getFieldJson(field Field) string {
	return `json:"` + field.Field + `"`
}

// 获取字段json描述
func getFieldJson2(field Field) string {

	// 是否主键
	keyJson := ``
	if field.Key == "PRI" {
		keyJson = ` pk`
	}

	// extra
	if field.Extra == "auto_increment" {
		keyJson += ` autoincr`
	}

	// 是否允许为Null
	notNullJson := ``
	if field.Null == "NO" {
		notNullJson = ` NOT NULL`
	}

	// 默认值
	defaultJson := ``
	if len(field.Default) > 0 {
		fieldType := getFiledType(field)
		if fieldType == "string" {
			defaultJson = ` default '` + field.Default + `'`
		} else {
			defaultJson = ` default ` + field.Default + ``
		}
	}

	// 字段备注
	commentJson := ``
	if len(field.Comment) > 0 {
		field.Comment = strings.Replace(field.Comment, "\n", "", -1)
		field.Comment = strings.Replace(field.Comment, "\r", "", -1)
		commentJson = ` comment('` + field.Comment + `')`
	}

	// 字段类型
	typeJson := ` ` + field.Type

	json := `json:"` + field.Field + `" xorm:"` + keyJson + notNullJson + defaultJson + commentJson + typeJson + `"`
	return json
}

// 获取字段说明
func getFieldComment(field Field) string {
	if len(field.Comment) > 0 {
		return "// " + field.Comment
	}
	return ""
}

// 检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

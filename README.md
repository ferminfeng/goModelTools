#### GO语言连接Mysql生成对应的model，包括对应字段类型、注释等。生成基础的结构体

**目前暂近支持xorm 生成示例---------**

```go 
package models

import "time"

func (dao Introducer) TableName() string {
	return "introducer"
}

// Introducer 介绍人表
type Introducer struct {
	Id                int       `json:"id" xorm:" pk NOT NULL auto_increment comment('id') int(11)"`
	Mobile            string    `json:"mobile" xorm:" NOT NULL comment('介绍人手机号') varchar(16)"`
	SystemType        int       `json:"system_type" xorm:" NOT NULL default 1 comment('介绍人系统类型') tinyint(1)"`
	Name              string    `json:"name" xorm:" NOT NULL comment('介绍人姓名') varchar(32)"`
	InterestingLifeId int       `json:"interesting_life_id" xorm:" NOT NULL comment('介绍人有趣生活ID') int(11)"`
	Remark            string    `json:"remark" xorm:" NOT NULL comment('备注') varchar(255)"`
	Status            int       `json:"status" xorm:" NOT NULL default 1 comment('状态 1-启用 2-禁用') tinyint(1)"`
	CreatedAt         time.Time `json:"created_at" xorm:" comment('创建时间') timestamp"`
	UpdatedAt         time.Time `json:"updated_at" xorm:" comment('更新时间') timestamp"`
	DeletedAt         time.Time `json:"deleted_at" xorm:" comment('删除时间') timestamp"`
}

```

**参数配置**

```
cp config.toml.template config.toml
修改config.toml内数据库配置

生成model
go run main.go 
    --c config.toml
    --model_path ./models/ 
    --model_replace true
```

| 参数名 | 必选   | 说明                                        |
|:----|:-----|-------------------------------------------|
| c   | 否    | 配置文件，默认为config.toml                       |
| model_path   | 否    | model文件生成位置，默认为 ./models/base,注意需要提前生成文件夹 |
| model_replace   | 否    | 是否覆盖已存在model，默认为true                      |
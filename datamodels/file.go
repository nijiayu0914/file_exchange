package datamodels

import "gorm.io/gorm"

// file 数据库file表映射实体
type File struct {
	gorm.Model // gorm 基本数据表字段
	Uuid string `gorm:"uniqueIndex; size:50" json:"uuid"` // 文件对象uuid
	FileName string `gorm:"index; size:50" json:"file_name"` // 文件对象名
	Capacity float64 `gorm:"default:0" json:"capacity"` // 文件存储最大容量限制
	UsageCapacity float64 `gorm:"default:0" json:"usage_capacity"` // 已用容量
	UserName string `gorm:"size:50" json:"user_name"` // 文件夹归属用户
	User     User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserName; references:UserName;"` // 对应用户
}
package datamodels

import "gorm.io/gorm"

// UserPlugin 数据库user_plugin表映射实体
type UserPlugin struct {
	gorm.Model // gorm 基本数据表字段
	Permission int16 `gorm:"default:1001" json:"permission"` // 权限
	MaxLibrary int `gorm:"default:3" json:"max_library"` // 最大library数量
	UserName string `gorm:"size:50;uniqueIndex" json:"user_name"` // 用户配置归属用户
	User     User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserName; references:UserName;"` // 对应用户
}

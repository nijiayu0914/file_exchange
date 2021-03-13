package datamodels

import "gorm.io/gorm"

// User 数据库user表映射实体
type User struct {
	gorm.Model // gorm 基本数据表字段
	UserName string `gorm:"uniqueIndex; size:50" json:"user_name"` // 用户名
	Password string `json:"password"` // 密码
}
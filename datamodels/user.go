package datamodels

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"uniqueIndex; size:50" json:"user_name"`
	Password string `json:"password"`
}
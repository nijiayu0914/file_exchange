package datamodels

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Uuid string `gorm:"uniqueIndex; size:50" json:"uuid"`
	FileName string `gorm:"index; size:50" json:"file_name"`
	Capacity float64 `gorm:"default:0" json:"capacity"`
	UsageCapacity float64 `gorm:"default:0" json:"usage_capacity"`
	UserName string `gorm:"size:50" json:"user_name"`
	User     User `gorm:"foreignKey:UserName; references:UserName;"`
}
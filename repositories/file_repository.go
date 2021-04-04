package repositories

import (
	"errors"
	"file_exchange/datamodels"
	"file_exchange/utils"
	"gorm.io/gorm"
)

// IFileRepository file表CRUD接口
type IFileRepository interface {
	SelectByUserName(userName string) (
		files[] map[string]interface{}, err error)
	SelectByFileUuid(fileUuid string) (file *datamodels.File, err error)
	SelectPaginate(page int, pageSize int, keyWord string) (
		files []datamodels.File, err error)
	Count(keyWord string) (count int64, err error)
	Insert(file *datamodels.File) (fileId uint, fileUuid string, err error)
	UpdateFileName(fileName string, uuid string) (err error)
	UpdateUsageCapacity(usageCapacity float64, uuid string) (err error)
	UpdateCapacity(capacity float64, uuid string) (err error)
	DeleteByUuid(fileUuid string)(err error)
}

// NewFileRepository 初始化file表操作对象
func NewFileRepository(db *gorm.DB) IFileRepository {
	return &FileManagerRepository{db}
}

// FileManagerRepository file表操作对象
type FileManagerRepository struct {
	db *gorm.DB // gorm Db
}

// SelectByUserName 查询用户下的所有文件夹
func (f *FileManagerRepository) SelectByUserName(userName string) (
	files[] map[string]interface{}, err error) {
	result := f.db.Model(&datamodels.File{}).Where(
		"user_name = ?", userName).Find(&files)
	if result.Error != nil{
		return files, result.Error
	}
	return files, nil
}

// SelectByFileUuid 根据uuid查询文件夹信息
func (f *FileManagerRepository) SelectByFileUuid(fileUuid string) (
	table *datamodels.File, err error) {
	var file datamodels.File
	result := f.db.Model(&datamodels.File{}).Where(
		"uuid = ?", fileUuid).First(&file)
	if result.Error != nil{
		return &file, result.Error
	}
	return &file, nil
}

// SelectPaginate 分页查询用户Library
func (f *FileManagerRepository) SelectPaginate(
	page int, pageSize int, keyWord string) (
	files []datamodels.File, err error) {
	result := f.db.Scopes(utils.Paginate(page, pageSize)).Where(
		"user_name LIKE ?", "%" + keyWord + "%").Or(
		"file_name LIKE ?", "%" + keyWord + "%").Or(
		"uuid LIKE ?", keyWord + "%").Find(&files)
	if result.Error != nil{
		return files, result.Error
	}
	return files, nil
}

// Count file表计数
func (f *FileManagerRepository) Count(keyWord string) (count int64, err error) {
	result := f.db.Model(&datamodels.File{}).Where(
		"user_name LIKE ?", "%" + keyWord + "%").Or(
		"file_name LIKE ?", "%" + keyWord + "%").Or(
		"uuid LIKE ?", keyWord + "%").Count(&count)
	if result.Error != nil{
		return 0, result.Error
	}
	return count, nil
}

// Insert 新增文件夹
func (f *FileManagerRepository) Insert(file *datamodels.File) (
	fileId uint, fileUuid string, err error) {
	result := f.db.Create(file)
	if result.Error != nil{
		return file.ID, file.Uuid, errors.New("创建文件夹失败")
	}
	return file.ID, file.Uuid, nil
}

// UpdateFileName 更新文件夹名称
func (f *FileManagerRepository) UpdateFileName(
	fileName string, uuid string) (err error) {
	result := f.db.Model(
		&datamodels.File{}).Where(
			"uuid = ?", uuid).Update("file_name", fileName)
	if result.Error != nil{
		return result.Error
	}
	return nil
}

// DeleteByUuid 根据uuid删除文件夹
func (f *FileManagerRepository) DeleteByUuid(fileUuid string) (err error) {
	result := f.db.Where("uuid = ?", fileUuid).Delete(&datamodels.File{})
	if result.Error != nil{
		return result.Error
	}
	return nil
}

func (f *FileManagerRepository) UpdateUsageCapacity(
	usageCapacity float64, uuid string) (err error) {
	result := f.db.Model(
		&datamodels.File{}).Where("uuid = ?", uuid).Update(
			"usage_capacity", usageCapacity)
	if result.Error != nil{
		return result.Error
	}
	return nil
}

func (f *FileManagerRepository) UpdateCapacity(
	capacity float64, uuid string) (err error) {
	result := f.db.Model(
		&datamodels.File{}).Where("uuid = ?", uuid).Update(
		"capacity", capacity)
	if result.Error != nil{
		return result.Error
	}
	return nil
}

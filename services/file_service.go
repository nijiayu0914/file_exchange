package services

import (
	"errors"
	"file_exchange/datamodels"
	"file_exchange/repositories"
	"github.com/google/uuid"
)

// IFileService file服务接口
type IFileService interface {
	FindFilesByUserName(userName string) (files[] map[string]interface{}, err error)
	FindFileByUuid(fileUuid string) (file *datamodels.File, err error)
	FindByPaginate(page int, pageSize int, keyWord string) (
		files []datamodels.File, err error)
	Count(keyWord string) (count int64, err error)
	CreateFile(file *datamodels.File) (fileId uint, fileUuid string, err error)
	UpdateFileName(fileName string, uuid string) (err error)
	UpdateUsageCapacity(usageCapacity float64, uuid string, how string) (err error)
	UpdateCapacity(capacity float64, uuid string) (err error)
	DeleteByUuid(fileUuid string) (err error)
	CheckCapacity(uuid string) (usageCapacity float64,
		capacity float64, free float64, err error)
}

// NewFileService 初始化file服务操作对象
func NewFileService(repository repositories.IFileRepository) IFileService {
	return &FileService{repository}
}

// FileService file服务操作对象
type FileService struct {
	FileRepository repositories.IFileRepository
}

// FindByPaginate 分页查询用户配置
func (f *FileService) FindByPaginate(page int, pageSize int, keyWord string) (
	files []datamodels.File, err error) {
	return  f.FileRepository.SelectPaginate(page, pageSize, keyWord)
}

// Count file表计数
func (f *FileService) Count(keyWord string) (count int64, err error) {
	return f.FileRepository.Count(keyWord)
}

// CreateFile 创建文件夹
func (f *FileService) CreateFile(file *datamodels.File) (fileId uint,
	fileUuid string, err error) {
	uuidObj, err := uuid.NewRandom()
	if err != nil{
		return 0, "", err
	}
	file.Uuid = uuidObj.String()
	return f.FileRepository.Insert(file)
}

// FindFilesByUserName 根据用户名查询文件
func (f *FileService) FindFilesByUserName(userName string) (
	files[] map[string]interface{}, err error) {
	return f.FileRepository.SelectByUserName(userName)
}

// FindFileByUuid 根据uuid查询文件
func (f *FileService) FindFileByUuid(fileUuid string) (
	file *datamodels.File, err error) {
	return f.FileRepository.SelectByFileUuid(fileUuid)
}

// UpdateFileName 修改文件夹名称
func (f *FileService) UpdateFileName(fileName string, uuid string) (err error) {
	return f.FileRepository.UpdateFileName(fileName, uuid)
}

// DeleteByUuid 根据uuid删除文件夹
func (f *FileService) DeleteByUuid(fileUuid string) (err error) {
	return f.FileRepository.DeleteByUuid(fileUuid)
}

// UpdateUsageCapacity 更新使用用量
func (f *FileService) UpdateUsageCapacity(usageCapacity float64,
	uuid string, how string) (err error) {
	file, err := f.FindFileByUuid(uuid)
	if err != nil{
		return err
	}
	if how == "increase"{
		usageCapacity = file.UsageCapacity + usageCapacity
	}else if how == "decrease"{
		usageCapacity = file.UsageCapacity - usageCapacity
	}else if how == "overwrite"{

	}else{
		return errors.New("写入方式错误")
	}
	return f.FileRepository.UpdateUsageCapacity(usageCapacity, uuid)
}

// UpdateCapacity 更新允许容量
func (f *FileService) UpdateCapacity(capacity float64,
	uuid string) (err error) {
	return f.FileRepository.UpdateCapacity(capacity, uuid)
}

// CheckCapacity 检查用量
func (f *FileService) CheckCapacity(uuid string) (usageCapacity float64,
	capacity float64, free float64, err error) {
	file, err := f.FindFileByUuid(uuid)
	if err != nil{
		return 0, 0, 0, err
	}
	return file.UsageCapacity, file.Capacity,
	file.Capacity - file.UsageCapacity, nil
}

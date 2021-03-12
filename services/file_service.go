package services

import (
	"errors"
	"file_exchange/datamodels"
	"file_exchange/repositories"
	"github.com/google/uuid"
)

type IFileService interface {
	FindFilesByUserName(userName string) (files[] map[string]interface{}, err error)
	FindFileByUuid(fileUuid string) (file *datamodels.File, err error)
	CreateFile(file *datamodels.File) (fileId uint, fileUuid string, err error)
	UpdateFileName(fileName string, uuid string) (err error)
	UpdateUsageCapacity(usageCapacity float64, uuid string, how string) (err error)
	UpdateCapacity(capacity float64, uuid string) (err error)
	DeleteByUuid(fileUuid string) (err error)
	CheckCapacity(uuid string) (usageCapacity float64,
		capacity float64, free float64, err error)
}

func NewFileService(repository repositories.IFileRepository) IFileService {
	return &FileService{repository}
}

type FileService struct {
	FileRepository repositories.IFileRepository
}

func (f *FileService) CreateFile(file *datamodels.File) (fileId uint,
	fileUuid string, err error) {
	uuidObj, err := uuid.NewRandom()
	if err != nil{
		return 0, "", err
	}
	file.Uuid = uuidObj.String()
	return f.FileRepository.Insert(file)
}

func (f *FileService) FindFilesByUserName(userName string) (
	files[] map[string]interface{}, err error) {
	return f.FileRepository.SelectByUserName(userName)
}

func (f *FileService) FindFileByUuid(fileUuid string) (
	file *datamodels.File, err error) {
	return f.FileRepository.SelectByFileUuid(fileUuid)
}

func (f *FileService) UpdateFileName(fileName string, uuid string) (err error) {
	return f.FileRepository.UpdateFileName(fileName, uuid)
}

func (f *FileService) DeleteByUuid(fileUuid string) (err error) {
	return f.FileRepository.DeleteByUuid(fileUuid)
}

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

func (f *FileService) UpdateCapacity(capacity float64,
	uuid string) (err error) {
	return f.FileRepository.UpdateCapacity(capacity, uuid)
}

func (f *FileService) CheckCapacity(uuid string) (usageCapacity float64,
	capacity float64, free float64, err error) {
	file, err := f.FindFileByUuid(uuid)
	if err != nil{
		return 0, 0, 0, err
	}
	return file.UsageCapacity, file.Capacity,
	file.Capacity - file.UsageCapacity, nil
}
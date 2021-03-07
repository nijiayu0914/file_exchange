package services

import (
	"file_exchange/datamodels"
	"file_exchange/repositories"
	"github.com/google/uuid"
)

type IFileService interface {
	FindFilesByUserName(userName string) (files[] map[string]interface{}, err error)
	CreateFile(file *datamodels.File) (fileId uint, fileUuid string, err error)
	UpdateFileName(fileName string, uuid string) (err error)
	DeleteByUuid(fileUuid string) (err error)
}

func NewFileService(repository repositories.IFileRepository) IFileService {
	return &FileService{repository}
}

type FileService struct {
	FileRepository repositories.IFileRepository
}

func (f *FileService) CreateFile(file *datamodels.File) (fileId uint, fileUuid string, err error) {
	uuidObj, err := uuid.NewRandom()
	if err != nil{
		return 0, "", err
	}
	file.Uuid = uuidObj.String()
	return f.FileRepository.Insert(file)
}

func (f *FileService) FindFilesByUserName(userName string) (files[] map[string]interface{}, err error) {
	return f.FileRepository.SelectByUserName(userName)
}

func (f *FileService) UpdateFileName(fileName string, uuid string) (err error) {
	return f.FileRepository.UpdateFileName(fileName, uuid)
}

func (f *FileService) DeleteByUuid(fileUuid string) (err error) {
	return f.FileRepository.DeleteByUuid(fileUuid)
}

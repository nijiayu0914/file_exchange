package services

import (
	"file_exchange/datamodels"
	"file_exchange/repositories"
)

// IUserPluginService UserPlugin服务接口
type IUserPluginService interface {
	FindByUserName(userName string) (userPlugin *datamodels.UserPlugin, err error)
	FindByPaginate(page int, pageSize int) (
		userPlugins []datamodels.UserPlugin, err error)
	Count() (count int64, err error)
	UpdateMaxLibrary(userName string, maxLibrary int) (err error)
	UpdatePermission(userName string, permission int16) (err error)
}

// NewUserPluginService 初始化UserPlugin服务操作对象
func NewUserPluginService(repository repositories.IUserPluginRepository) IUserPluginService {
	return &UserPluginService{repository}
}

// UserPluginService UserPlugin服务操作对象
type UserPluginService struct {
	UserPluginRepository repositories.IUserPluginRepository
}

// Find 根据用户名查询用户配置
func (u *UserPluginService) FindByUserName(userName string) (
	userPlugin *datamodels.UserPlugin, err error) {
	return u.UserPluginRepository.Select(userName)
}

// FindByUserNamePaginate 分页查询用户配置
func (u *UserPluginService) FindByPaginate(page int, pageSize int) (
	userPlugins []datamodels.UserPlugin, err error) {
	return  u.UserPluginRepository.SelectPaginate(page, pageSize)
}

// Count user_plugin表计数
func (u *UserPluginService) Count() (count int64, err error) {
	return u.UserPluginRepository.Count()
}

// UpdateMaxLibrary 更新用户library最大数量
func (u *UserPluginService) UpdateMaxLibrary(
	userName string, maxLibrary int) (err error) {
	return u.UserPluginRepository.UpdateMaxLibrary(userName, maxLibrary)
}

// UpdatePermission 更新用户permission等级
func (u *UserPluginService) UpdatePermission(
	userName string, permission int16) (err error) {
	return u.UserPluginRepository.UpdatePermission(userName, permission)
}

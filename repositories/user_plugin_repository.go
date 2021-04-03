package repositories

import (
	"errors"
	"file_exchange/datamodels"
	"file_exchange/utils"
	"gorm.io/gorm"
)

// IUserPluginRepository user_plugin表CRUD接口
type IUserPluginRepository interface {
	Insert(user *datamodels.UserPlugin) (userPluginId uint, err error)
	Select(userName string) (userPlugin *datamodels.UserPlugin, err error)
	SelectPaginate(page int, pageSize int) (
		userPlugins []datamodels.UserPlugin, err error)
	Count() (count int64, err error)
	UpdateMaxLibrary(userName string, maxLibrary int) (err error)
	UpdatePermission(userName string, permission int16) (err error)
}

// NewUserPluginRepository 初始化user_plugin表操作对象
func NewUserPluginRepository(db *gorm.DB) IUserPluginRepository {
	return &UserPluginManagerRepository{db}
}

// UserPluginManagerRepository user_plugin表操作对象
type UserPluginManagerRepository struct {
	db *gorm.DB // gorm Db
}

// Insert 根据用户名生成用户配置
func (u *UserPluginManagerRepository) Insert(
	userPlugin *datamodels.UserPlugin) (userPluginId uint, err error) {
		result := u.db.Create(userPlugin)
		if result.Error != nil{
			return userPlugin.ID, errors.New("创建用户配置失败")
		}
		return userPlugin.ID, nil
}

// Select 根据用户名查询用户配置
func (u *UserPluginManagerRepository) Select(userName string) (
	table *datamodels.UserPlugin, err error) {
	var userPlugin datamodels.UserPlugin
	result := u.db.Where("user_name = ?", userName).First(&userPlugin)
	if result.Error != nil{
		return &userPlugin, result.Error
	}
	return &userPlugin, nil
}

// Count user_plugin表计数
func (u *UserPluginManagerRepository) Count() (count int64, err error) {
	result := u.db.Model(&datamodels.UserPlugin{}).Count(&count)
	if result.Error != nil{
		return 0, result.Error
	}
	return count, nil
}

// SelectByUserNamePaginate 分页查询用户配置
func (u *UserPluginManagerRepository) SelectPaginate(
	page int, pageSize int) (
	userPlugins []datamodels.UserPlugin, err error) {
	result := u.db.Scopes(utils.Paginate(page, pageSize)).Find(&userPlugins)
	if result.Error != nil{
		return userPlugins, result.Error
	}
	return userPlugins, nil
}

// UpdateMaxLibrary 更新用户library最大数量
func (u *UserPluginManagerRepository) UpdateMaxLibrary(
	userName string, maxLibrary int) (err error) {
	result := u.db.Model(&datamodels.UserPlugin{}).Where(
		"User_name = ?", userName).Update("max_library", maxLibrary)
	if result.Error != nil{
		return result.Error
	}
	return nil
}

// UpdatePermission 更新用户permission等级
func (u *UserPluginManagerRepository) UpdatePermission(
	userName string, permission int16) (err error) {
	result := u.db.Model(&datamodels.UserPlugin{}).Where(
		"User_name = ?", userName).Update("permission", permission)
	if result.Error != nil{
		return result.Error
	}
	return nil
}


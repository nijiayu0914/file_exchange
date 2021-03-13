// repositories 数据库操作层
package repositories

import (
	"errors"
	"file_exchange/datamodels"
	"gorm.io/gorm"
)

// IUserRepository user表CRUD接口
type IUserRepository interface {
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId uint, err error)
	Update(user *datamodels.User) (userId uint, err error)
	SelectAll() (users []datamodels.User, err error)

}

// NewUserRepository 初始化user表操作对象
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserManagerRepository{db}
}

// UserManagerRepository user表操作对象
type UserManagerRepository struct {
	db *gorm.DB // gorm Db
}

// Select 根据用户名查询用户
func (u *UserManagerRepository) Select(userName string) (
	table *datamodels.User, err error) {
	var user datamodels.User
	result := u.db.Where("user_name = ?", userName).First(&user)
	if result.Error != nil{
		return &user, result.Error
	}
	return &user, nil
}

// Insert 新增用户
func (u *UserManagerRepository) Insert(user *datamodels.User) (
	tableId uint, err error) {
	result := u.db.Create(user)
	if result.Error != nil{
		return user.ID, errors.New("用户名重复")
	}
	return user.ID, nil
}

// SelectAll 查询所有用户
func (u *UserManagerRepository) SelectAll() (
	users []datamodels.User, err error) {
	result := u.db.Find(&users)
	if result.Error != nil{
		return users, result.Error
	}
	return users, nil
}

func (u *UserManagerRepository) Update(user *datamodels.User) (userId uint, err error) {
	result := u.db.Save(user)
	if result.Error != nil{
		return user.ID, result.Error
	}
	return user.ID, nil
}

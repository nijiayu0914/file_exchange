// services 服务层
package services

import (
	"file_exchange/datamodels"
	"file_exchange/repositories"
	"golang.org/x/crypto/bcrypt"
)

// IUserService user服务接口
type IUserService interface {
	AddUser(user *datamodels.User) (userId uint, err error)
	FindUser(userName string) (user *datamodels.User, err error)
	FindAllUser() (users [] datamodels.User, err error)
	ChangePassword(user *datamodels.User, newPassword string) (userId uint, err error)
}

// NewUserService 初始化user服务操作对象
func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

// UserService user服务操作对象
type UserService struct {
	UserRepository repositories.IUserRepository
}

// FindUser 根据用户名查询用户信息
func (u *UserService) FindUser(userName string) (user *datamodels.User, err error) {
	return u.UserRepository.Select(userName)
}

// AddUser 新增用户
func (u *UserService) AddUser(user *datamodels.User) (userId uint, err error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.ID, err
	}
	user.Password = string(hashPassword)
	return u.UserRepository.Insert(user)
}

// FindAllUser 查询所有用户
func (u *UserService) FindAllUser() (users []datamodels.User, err error) {
	return u.UserRepository.SelectAll()
}

// ChangePassword 修改密码
func (u *UserService) ChangePassword(user *datamodels.User, newPassword string) (userId uint, err error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return user.ID, err
	}
	user.Password = string(hashPassword)
	return u.UserRepository.Update(user)
}

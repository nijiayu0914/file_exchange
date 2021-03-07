package services

import (
	"file_exchange/datamodels"
	"file_exchange/repositories"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	AddUser(user *datamodels.User) (userId uint, err error)
	FindUser(userName string) (user *datamodels.User, err error)
	FindAllUser() (users [] datamodels.User, err error)
	ChangePassword(user *datamodels.User, newPassword string) (userId uint, err error)
}

func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (u *UserService) FindUser(userName string) (user *datamodels.User, err error) {
	return u.UserRepository.Select(userName)
}

func (u *UserService) AddUser(user *datamodels.User) (userId uint, err error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.ID, err
	}
	user.Password = string(hashPassword)
	return u.UserRepository.Insert(user)
	}

func (u *UserService) FindAllUser() (users []datamodels.User, err error) {
	return u.UserRepository.SelectAll()
}

func (u *UserService) ChangePassword(user *datamodels.User, newPassword string) (userId uint, err error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return user.ID, err
	}
	user.Password = string(hashPassword)
	return u.UserRepository.Update(user)
}
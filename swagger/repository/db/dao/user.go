package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"swagger/repository/db/model"
	req "swagger/request"
)

type UserDao struct {
	*gorm.DB
}

// NewUserDao 返回 db 的链接
func NewUserDao(context context.Context) *UserDao {
	return &UserDao{NewDBClient(context)}
}

// UserList 用户列表
func (dao *UserDao) UserList() []model.User {
	var users []model.User
	dao.Model(&model.User{}).Find(&users)
	return users
}

// GetUserInfoById 根据用户ID获取信息
func (dao *UserDao) GetUserInfoById(userId int64) (r *model.User, err error) {
	err = dao.Model(&model.User{}).Where("id=?", userId).First(&r).Error
	return
}

// GetUserInfo 根据用户名获取用户信息
func (dao *UserDao) GetUserInfo(req req.LoginRequest) (r *model.User, err error) {
	err = dao.Model(&model.User{}).Where("username=?", req.Username).First(&r).Error
	return
}

// CreateUser 创建用户
func (dao *UserDao) CreateUser(req req.RegisterRequest) (err error) {
	var user model.User
	var count int64
	dao.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)

	if count != 0 {
		return errors.New("username Exist")
	}

	user = model.User{
		Username: req.Username,
		Password: req.Password,
		Mobile:   req.Mobile,
	}

	if err = dao.Model(&model.User{}).Create(&user).Error; err != nil {
		return
	}
	return nil
}

// UpdateUser 修改用户
func (dao *UserDao) UpdateUser(id int64, req req.RegisterRequest) (err error) {
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Mobile:   req.Mobile,
	}
	err = dao.Model(&model.User{}).Where("id = ?", id).Updates(&user).Error
	return err
}

// DeleteUserById 根据用户ID删除用户
func (dao *UserDao) DeleteUserById(id int64) (err error) {
	err = dao.Model(&model.User{}).
		Where("id = ?", id).
		Delete(model.User{}).Error
	return err
}

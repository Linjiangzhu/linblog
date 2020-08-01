package repository

import (
	"github.com/Linjiangzhu/linblog/linblog-backend/model"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) CreateRole(r *model.Role) error {
	return ur.db.Create(r).Error
}

func (ur *UserRepository) DeleteRole(r *model.Role) error {
	return ur.db.Delete(r).Error
}

func (ur *UserRepository) UpdateRole(r *model.Role) error {
	return ur.db.Model(r).Update("name", r.Name).Error
}

func (ur *UserRepository) GetRoles() (*[]model.Role, error) {
	var roles []model.Role
	err := ur.db.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

func (ur *UserRepository) CreateUser(u *model.User) error {
	return ur.db.Create(u).Error
}

func (ur *UserRepository) DeleteUser(u *model.User) error {
	return ur.db.Delete(u).Error
}

func (ur *UserRepository) UpdateUser(u *model.User) error {
	return ur.db.Model(u).Updates(map[string]interface{}{
		"user_name":  u.UserName,
		"email":      u.Email,
		"nick_name":  u.NickName,
		"password":   u.Password,
		"role_id":    u.RoleID,
		"last_login": u.LastLogin,
	}).Error
}

func (ur *UserRepository) GetUserByID(id string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByUserName(username string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUsers() (*[]model.User, error) {
	var users []model.User
	err := ur.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return &users, nil
}

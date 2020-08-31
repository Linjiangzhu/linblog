package repository

import (
	"github.com/Linjiangzhu/blog-v2/model"
)

const (
	getUserByIDSQL = "SELECT * from users where id = ?"
	insertUserSQL  = "INSERT INTO users (id, username, password, nickname) VALUES(?, ?, ?, ?)"
)

func (r *Repository) GetUserByID(uid string) (*model.User, error) {
	u := model.User{}
	if err := r.db.DB().QueryRow(getUserByIDSQL, uid).Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Username,
		&u.Password,
		&u.NickName,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) CreateUser(u *model.User) (*model.User, error) {
	_, err := r.db.DB().Exec(insertUserSQL, u.ID, u.Username, u.Password, u.NickName)
	if err != nil {
		return nil, err
	}
	return r.GetUserByID(u.ID)
}

func (r *Repository) GetUserByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

package repository

import "github.com/Linjiangzhu/linblog/linblog-backend/model"

const (
	getUserByIDSQL = "SELECT * from users where id = ?"
	insertUserSQL  = "INSERT INTO users (id, username, password, nickname) VALUES(?, ?, ?, ?)"
)

func (r *Repository) GetUserByID(uid string) (*model.User, error) {
	u := model.User{}
	sqlDB, _ := r.db.DB()
	if err := sqlDB.QueryRow(getUserByIDSQL, uid).Scan(
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
	sqlDB, _ := r.db.DB()
	_, err := sqlDB.Exec(insertUserSQL, u.ID, u.Username, u.Password, u.NickName)
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

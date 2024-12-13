package mysqluser

import (
	"database/sql"
	"event-manager/entity"
	"event-manager/repository/mysql"
	"fmt"
)

type UserRepo struct {
	conn *mysql.MySQLDB
}

func New(c *mysql.MySQLDB) UserRepo {
	return UserRepo{
		conn: c,
	}
}

type userModel struct {
	id             uint
	username       string
	hashedPassword string
}

func (e *userModel) ToUserEntity() entity.User {
	var entiy entity.User

	entiy.SetID(e.id)
	entiy.SetUsername(e.username)
	entiy.SetPassword(e.hashedPassword)

	return entiy
}

func (r UserRepo) GetUserByUsername(username string) (entity.User, bool, error) {
	var model userModel

	row := r.conn.Conn().QueryRow(`select id, username, hashed_password from users where username = ?`, username)
	err := row.Scan(&model.id, &model.username, &model.hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}

		return entity.User{}, false, err
	}

	return model.ToUserEntity(), true, nil
}

func (r UserRepo) CreateUser(u entity.User) (entity.User, error) {
	res, err := r.conn.Conn().Exec(`insert into users (username, hashed_password) values (?, ?)`,
		u.Username(), u.Password())

	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command: %w", err)
	}

	// error is always nil
	id, _ := res.LastInsertId()
	u.SetID(uint(id))

	return u, nil
}

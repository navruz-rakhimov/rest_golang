package data

import (
	"database/sql"
	"errors"
)

type User struct {
	Id       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
}

type Users struct {
	DB *sql.DB
}

func (users Users) Insert(user *User) error {
	query := `INSERT INTO users (login, password, name, age) VALUES($1, $2, $3, $4) RETURNING id`
	args := []interface{}{user.Login, user.Password, user.Name, user.Age}
	return users.DB.QueryRow(query, args...).Scan(&user.Id)
}

func (users Users) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, login, password, name, age FROM users WHERE id = $1`

	var user User
	err := users.DB.QueryRow(query, id).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
		&user.Name,
		&user.Age,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (users Users) GetByName(name string) (*User, error) {
	query := `SELECT id, login, password, name, age FROM users WHERE name=$1`

	var user User
	err := users.DB.QueryRow(query, name).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
		&user.Name,
		&user.Age,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (users Users) GetByLogin(login string) (*User, error) {
	query := `SELECT id, login, password, name, age FROM users WHERE login = $1`

	var user User
	err := users.DB.QueryRow(query, login).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
		&user.Name,
		&user.Age,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

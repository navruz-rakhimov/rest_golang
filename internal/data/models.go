package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users  Users
	Phones Phones
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: Users{
			DB: db,
		},
		Phones: Phones{
			DB: db,
		},
	}
}

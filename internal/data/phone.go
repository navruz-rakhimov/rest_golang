package data

import (
	"database/sql"
	"errors"
)

type Phone struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
	IsFax       bool   `json:"is_fax"`
}

type Phones struct {
	DB *sql.DB
}

func (phones Phones) Get(id int) (*Phone, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, user_id, phone, description, is_fax FROM phones WHERE id = $1`

	var phone Phone
	err := phones.DB.QueryRow(query, id).Scan(
		&phone.Id,
		&phone.UserId,
		&phone.Phone,
		&phone.Description,
		&phone.IsFax,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &phone, nil
}

func (phones Phones) Update(phone *Phone) error {
	query := `UPDATE phones 
			  SET 
			      user_id = $1,
			      phone = $2, 
			      description = $3, 
			      is_fax = $4
			  WHERE id = $5`

	args := []interface{}{
		phone.UserId,
		phone.Phone,
		phone.Description,
		phone.IsFax,
		phone.Id,
	}
	_, err := phones.DB.Exec(query, args...)
	return err
}

func (phones Phones) Insert(phone *Phone) error {
	query := `INSERT INTO phones (user_id, phone, description, is_fax) VALUES($1, $2, $3, $4) RETURNING id`
	args := []interface{}{phone.UserId, phone.Phone, phone.Description, phone.IsFax}
	return phones.DB.QueryRow(query, args...).Scan(&phone.Id)
}

func (phones Phones) GetAllWithNumber(phone string) ([]*Phone, error) {
	query := `SELECT id, user_id, phone, description, is_fax
			  FROM phones WHERE phone LIKE '%' || $1 || '%'`

	rows, err := phones.DB.Query(query, phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phoneNums []*Phone
	for rows.Next() {
		p := &Phone{}
		err = rows.Scan(
			&p.Id,
			&p.UserId,
			&p.Phone,
			&p.Description,
			&p.IsFax,
		)
		if err != nil {
			return nil, err
		}

		phoneNums = append(phoneNums, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return phoneNums, nil
}

func (phones Phones) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM phones WHERE id = $1`
	result, err := phones.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

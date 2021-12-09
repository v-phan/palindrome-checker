package db

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DbConnection struct {
	Db  *pgxpool.Pool
	Ctx context.Context
}

type User struct {
	Id        int32  `json: "id"`
	FirstName string `json: "firstName"`
	LastName  string `json: "lastName"`
}

func (u *User) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(u)
}

func (d *DbConnection) CreateUser(firstName string, lastName string) (*User, error) {
	sql := `INSERT INTO "user" (first_name, last_name) VALUES ($1, $2) RETURNING id`
	user := &User{FirstName: firstName, LastName: lastName}
	err := d.Db.QueryRow(d.Ctx, sql, firstName, lastName).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *DbConnection) GetUser(userId int32) (*User, error) {
	sql := `SELECT id, first_name, last_name FROM "user" WHERE id = $1`
	user := &User{}
	err := d.Db.QueryRow(d.Ctx, sql, userId).Scan(&user.Id, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *DbConnection) DeleteUser(userId int32) (bool, error) {
	sql := `DELETE FROM "user" WHERE id = $1`
	res, err := d.Db.Exec(d.Ctx, sql, userId)
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}

func (d *DbConnection) UpdateUser(userId int32, firstName string, lastName string) (*User, error) {
	user := &User{}

	sql := `UPDATE "user" SET first_name = $1, last_name = $2 WHERE id = $3 RETURNING *`

	err := d.Db.QueryRow(d.Ctx, sql, firstName, lastName, userId).Scan(&user.Id, &user.FirstName, &user.LastName)
	if err != nil {
		log.Println(sql)
		return nil, err
	}

	return user, nil
}

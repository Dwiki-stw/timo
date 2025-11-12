package models

import "time"

type User struct {
	ID        int64     `db:"id"`
	Uid       string    `db:"uid"`
	GoogleID  *string   `db:"google_id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  *string   `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

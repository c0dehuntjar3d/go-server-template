package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id        int64
	Uuid      string
	Login     string
	Password  string
	CreatedAt time.Time
	UpdatedAt *sql.NullTime
	DeletedAt *sql.NullTime
}

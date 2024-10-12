package domain

import (
	"time"
)

type User struct {
	Uuid      string
	Login     string
	Password  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

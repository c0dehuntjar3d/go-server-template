package domain

import (
	"time"
)

type User struct {
	Login     string
	Password  string
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type UserInfo struct {
	FirstName string
	LastName  string
	Age       int64
}

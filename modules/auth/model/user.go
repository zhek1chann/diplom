package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Info      *UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserInfo struct {
	Name        string
	PhoneNumber string
	Role        int
}

type AuthUser struct {
	ID             int64
	Info           *UserInfo
	Password       string
	HashedPassword string
}

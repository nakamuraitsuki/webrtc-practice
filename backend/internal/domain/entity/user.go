package entity

import (
	"time"
)
type User struct {
	id         int
	name       string
	email      string
	passwdhash string
	createdAt  time.Time
	updatedAt  *time.Time
}

func NewUser(id int, name, email, passwdhash string, createdAt time.Time, updatedAt *time.Time) *User {
	return &User{
		id:         id,
		name:       name,
		email:      email,
		passwdhash: passwdhash,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

func (u User) GetID() int {
	return u.id
}

func (u User) GetPasswdHash() string {
	return u.passwdhash
}

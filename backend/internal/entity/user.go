package entity

import "time"

type User struct {
	id         int        `json:"id"`
	name       string     `json:"name"`
	email      string     `json:"email"`
	passwdhash string     `json:"passwdhash"`
	createdAt  time.Time  `json:"created_at"`
	updatedAt  *time.Time `json:"updated_at"`
}

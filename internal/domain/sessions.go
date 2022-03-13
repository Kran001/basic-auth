package domain

import (
	"time"
)

type User struct {
	Id       int64  `json:"user_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Session struct {
	Id          int64     `json:"session_id"`
	Token       string    `json:"token"`
	SessionTime time.Time `json:"session_time"`
	User
}

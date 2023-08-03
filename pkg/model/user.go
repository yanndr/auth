package model

type User struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

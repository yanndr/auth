package models

type User struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	// Could have more fields firstname... but I focused on username/pwd
}

package services

import (
	"fmt"
)

type UsernameAlreadyExistErr struct {
	name string
}

func (e UsernameAlreadyExistErr) Error() string {
	return fmt.Sprintf("the entity with username:%s already exist", e.name)
}

var AutenticationErr = fmt.Errorf("authentication failed")

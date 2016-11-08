package models

import (
	// Native packages
	"time"

	// 3rd party packages
	"gopkg.in/guregu/null.v3"

	// Local packages
	"jaha-api/utils"
)

const USER_ROLE_GUEST = 1
const USER_ROLE_MOD = 2
const USER_ROLE_ADMIN = 3

type User struct {
	ID        int       `json:"-"`
	UUID      string    `json:"uuid" validate:"required,len=8"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-"`
	AuthKey   string    `json:"-" validate:"required,len=16`
	Role      int       `json:"-"`
	UpdatedAt null.Time `json:"updatedAt"`
	DeletedAt null.Time `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	errors    []string
}

type Users []User

type UserPayload struct {
	FirstName       string `json:"firstName" validate:"omitempty,gte=3"`
	LastName        string `json:"lastName" validate:"omitempty,gte=3"`
	Email           string `json:"email" validate:"omitempty,email"`
	Password        string `json:"password" validate:"omitempty,eqfield=PasswordConfirm"`
	PasswordConfirm string `json:"passwordConfirm" validate:"omitempty,gte=6"`
}

func (user *User) Valid() bool {
	validationError, validationErrors := utils.Validate(user)

	if validationError != nil {
		user.SetErrors(validationErrors)
		return false
	}

	return true
}

func (user *User) GetErrors() []string {
	return user.errors
}

func (user *User) SetErrors(errors []string) {
	user.errors = errors
}

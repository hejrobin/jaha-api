package models

import (
	// Native packages
	"time"

	// 3rd party packages
	"gopkg.in/guregu/null.v3"

	// Local packages
	"jaha-api/utils"
)

type Statement struct {
	ID         int       `json:"-"`
	UUID       string    `json:"uuid" validate:"required,len=8"`
	Body       string    `json:"body"`
	Category   Category  `json:"category"`
	CategoryId int       `json:"-"`
	UpdatedAt  null.Time `json:"updatedAt"`
	DeletedAt  null.Time `json:"-"`
	CreatedAt  time.Time `json:"createdAt"`
	errors     []string
}

type Statements []Statement

type StatementPayload struct {
	Body     string `json:"body" validate:"omitempty,gte=3"`
	Category string `json:"category" validate:"omitempty,len=8"`
}

func (statement *Statement) Valid() bool {
	validationError, validationErrors := utils.Validate(statement)

	if validationError != nil {
		statement.SetErrors(validationErrors)
		return false
	}

	return true
}

func (statement *Statement) GetErrors() []string {
	return statement.errors
}

func (statement *Statement) SetErrors(errors []string) {
	statement.errors = errors
}

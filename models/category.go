package models

import (
	// Native packages
	"time"

	// 3rd party packages
	"gopkg.in/guregu/null.v3"
)

type Category struct {
	ID        int       `json:"-"`
	UUID      string    `json:"uuid" validate:"required,len=8"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	UpdatedAt null.Time `json:"updatedAt"`
	DeletedAt null.Time `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type Categories []Category

type CategoryPayload struct {
	Name string `json:"name" validate:"omitempty,gte=3"`
	Slug string `json:"slug" validate:"omitempty,gte=3"`
}

func (category *Category) Valid() bool {
	validationError, validationErrors := utils.Validate(category)

	if validationError != nil {
		category.errors = validationErrors
		return false
	}

	return true
}

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

type CategoryCollection struct {
	Meta       Collection `json:"meta"`
	Collection Categories `json:"collection"`
}

type CategoryPayload struct {
	Name string `json:"name" validate:"omitempty,gte=3"`
	Slug string `json:"slug" validate:"omitempty,gte=3"`
}

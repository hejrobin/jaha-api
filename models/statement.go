package models

import (
	// Native packages
	"time"

	// 3rd party packages
	"gopkg.in/guregu/null.v3"
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
}

type Statements []Statement

type StatementCollection struct {
	Meta       Collection `json:"meta"`
	Collection Statements `json:"collection"`
}

type StatementPayload struct {
	Body     string `json:"body" validate:"omitempty,gte=3"`
	Category string `json:"category" validate:"omitempty,len=8"`
}

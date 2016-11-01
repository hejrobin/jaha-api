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
	Published int       `json:"published"`
	UpdatedAt null.Time `json:"updatedAt"`
	DeletedAt null.Time `json:"deletedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type Categories []Category

type CategoryCollection struct {
	Meta       Collection `json:"meta"`
	Collection Categories `json:"collection"`
}

type CategoryPayload struct {
	Name      string     `json:"name" validate:"omitempty,gte=3"`
	Slug      string     `json:"slug" validate:"omitempty,gte=3"`
	Published int        `json:"published" validate:"omitempty"`
	DeletedAt *time.Time `json:"deletedAt" validate:"omitempty"`
}

/**
 *	Pre destroy hook.
 *	Sets published state to false.
 *
 *	@param category Category - Category resource pointer.
 *
 *	@return error
 */
func (category *Category) BeforeDelete() (err error) {
	category.Published = 0
	return
}

package scopes

import (
	// 3rd party packages
	"github.com/jinzhu/gorm"
)

type categoryScopes struct{}

func (categoryScopes) Published(dbc *gorm.DB) *gorm.DB {
	return dbc.Where("`published` = 1")
}

func (categoryScopes) NotPublished(dbc *gorm.DB) *gorm.DB {
	return dbc.Where("`published` = 0")
}

func (categoryScopes) Deleted(dbc *gorm.DB) *gorm.DB {
	return dbc.Where("`deleted_at` IS NOT NULL")
}

func Category() categoryScopes {
	var scopes categoryScopes
	return scopes
}

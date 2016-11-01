package scopes

import (
	// 3rd party packages
	"github.com/jinzhu/gorm"
)

type statementScopes struct{}

func (statementScopes) Deleted(dbc *gorm.DB) *gorm.DB {
	return dbc.Where("`deleted_at` IS NOT NULL")
}

func Statement() statementScopes {
	var scopes statementScopes
	return scopes
}

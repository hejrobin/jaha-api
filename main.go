package main

import (
	// Local packages
	"jaha-api/db"
	"jaha-api/env"
	"jaha-api/routers"
)

func main() {
	dbc := db.GetConnection()

	dbc.SingularTable(true)

	if env.IsDevelopmentMode() {
		dbc.LogMode(true)
	}

	defaultRouter := routers.GetDefaultRouter()

	defaultRouter.Run(":" + env.GetPort())
}

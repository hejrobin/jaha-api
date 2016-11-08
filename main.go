package main

import (
	// Native packages
	"fmt"

	// 3rd party packages
	"github.com/gin-gonic/contrib/sessions"

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

	sessionStore := sessions.NewCookieStore([]byte(env.GetSessionKey()))
	sessionManager := sessions.Sessions(env.GetRealmKey(), sessionStore)

	defaultRouter := routers.GetDefaultRouter(sessionManager)

	fmt.Println(fmt.Sprintf("[%s] Initializing...", env.GetAppName()))

	defaultRouter.Run(":" + env.GetPort())
}

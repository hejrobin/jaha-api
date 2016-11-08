package middlewares

import (
	// Native packages
	"time"

	// 3rd party packages
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/appleboy/gin-jwt.v2"

	// Local packages
	"jaha-api/db"
	"jaha-api/env"
	"jaha-api/models"
	"jaha-api/responders"
	"jaha-api/utils"
)

/**
 *	Authentication middleware, validates user and user password.
 *
 *	@param userEmail string - User email, used as username.
 *	@param userPassword string - User password, unhashed.
 *	@param ctx *gin.Context - Gin context.
 *
 *	@return string, bool
 */
func AuthAuthenticator(userEmail string, userPassword string, ctx *gin.Context) (string, bool) {
	var user models.User

	db.GetConnection().Where("`email` = ?", userEmail).First(&user)
	passwordMatches := utils.PasswordMatch(user.Password, userPassword)

	ctx.Set("userId", user.UUID)

	return userEmail, passwordMatches
}

/**
 *	Authorizator middleware, validates whether or
 *
 *	@param userEmail string - User email, used as username.
 *	@param ctx *gin.Context - Gin context.
 *
 *	@return bool
 */
func AuthAuthorizator(userEmail string, ctx *gin.Context) bool {
	var user models.User

	db.GetConnection().Where("`email` = ?", userEmail).First(&user)

	session := sessions.Default(ctx)

	// @NOTE Set user ID even if it's nil.
	session.Set("userId", user.ID)
	session.Save()

	if user.AuthKey != "" {
		return true
	}

	return false
}

/**
 *	Unauthorized middleware, sends un authorized access JSON response.
 *
 *	@param ctx *gin.Context - Gin context.
 *	@param errorCode int - Authentication error code.
 *	@param errorMessage string - Authentication error message.
 *
 *	@return void
 */
func AuthUnauthorized(ctx *gin.Context, errorCode int, errorMessage string) {
	responders.Text().Unauthorized(ctx, errorMessage)
	return
}

/**
 *	Returns authentication middleware object.
 *
 *	@return *jwt.GinJWTMiddleware
 */
func AuthMiddleware() *jwt.GinJWTMiddleware {
	realmKey := env.GetRealmKey()

	middleware := &jwt.GinJWTMiddleware{
		Realm:         env.GetAppName(),
		Key:           []byte(realmKey),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		Authenticator: AuthAuthenticator,
		Authorizator:  AuthAuthorizator,
		Unauthorized:  AuthUnauthorized,
		TokenLookup:   "header:Authorization",
	}

	return middleware
}

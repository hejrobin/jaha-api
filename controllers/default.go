package controllers

import (
	// 3rd party packages
	"github.com/gin-gonic/gin"
)

/**
 *	Empty controller struct.
 */
type controllerPrototype struct{}

/**
 *	Missing route error handler.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (controllerPrototype) MissingRoute(ctx *gin.Context) {
	ctx.JSON(404, gin.H{
		"error": "Missing route handler.",
	})
}

/**
 *	Root route handler.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (controllerPrototype) Index(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Hello!",
	})
}

/**
 *	Returns instanciated "controller".
 *	@NOTE Classes aren't present in Go, return a struct with field methods instead.
 *
 *	@return controllerPrototype
 */
func DefaultController() controllerPrototype {
	var defaultController controllerPrototype
	return defaultController
}

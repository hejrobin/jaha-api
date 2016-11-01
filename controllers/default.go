package controllers

import (
	// 3rd party packages
	"github.com/gin-gonic/gin"
)

type defaultPrototype struct{}

/**
 *	Missing route error handler.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (defaultPrototype) MissingRoute(ctx *gin.Context) {
	ctx.JSON(404, gin.H{
		"error": "Missing route handler.",
	})
}

/**
 *	Returns instanciated "controller".
 *	@NOTE Classes aren't present in Go, return a struct with field methods instead.
 *
 *	@return defaultPrototype
 */
func DefaultController() defaultPrototype {
	var controllerInstance defaultPrototype
	return controllerInstance
}

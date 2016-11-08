package controllers

import (
	// 3rd party packages
	"github.com/gin-gonic/gin"

	// Local packages
	"jaha-api/responders"
)

const COLLECTION_DEFAULT_LIMIT = 25

type defaultPrototype struct{}

/**
 *	Missing route error handler.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (defaultPrototype) MissingRoute(ctx *gin.Context) {
	responders.Text().NotFound(ctx, "Route not found.")
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

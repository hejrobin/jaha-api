package constraints

import (
	// 3rd party packages
	"github.com/gin-gonic/gin"
)

type Constraint struct {
	RequestPath   string
	RequestMethod string
	Guard         func(ctx *gin.Context) bool
}

type Constraints []Constraint

var registeredConstraints Constraints

/**
 *	Adds a route constraint using request method, path and a guard handler.
 *
 *	@param requestMethod string
 *	@param requestPath string
 *	@param guardHandler func(ctx *gin.Context) bool
 *
 *	@return void
 */
func AddConstraint(requestMethod string, requestPath string, guardHandler func(ctx *gin.Context) bool) {
	registeredConstraints = append(registeredConstraints, Constraint{
		RequestPath:   requestPath,
		RequestMethod: requestMethod,
		Guard:         guardHandler,
	})
}

/**
 *	Returns all registered constraints.
 *
 *	@return Constraints
 */
func GetConstraints() Constraints {
	return registeredConstraints
}

package middlewares

import (
	// 3rd party packages
	"github.com/gin-gonic/gin"

	// Local packages
	"jaha-api/constraints"
)

/**
 *	Loops through registered constraints (if any) and validates them.
 *
 *	@return gin.HandlerFunc
 */
func Constraints() gin.HandlerFunc {
	// @NOTE Manually invoke constraints register functions
	constraints.UserConstraints()

	return func(ctx *gin.Context) {
		canContinueRequest := true
		registeredConstraints := constraints.GetConstraints()

		if len(registeredConstraints) > 0 {
			for _, constraint := range registeredConstraints {
				if constraint.RequestMethod == ctx.Request.Method && constraint.RequestPath == ctx.Request.URL.Path {
					if !constraint.Guard(ctx) {
						canContinueRequest = false
						break
					}
				}
			}

			if !canContinueRequest {
				ctx.JSON(403, gin.H{
					"error": "Permission denied.",
				})
				return
			}
		}

		ctx.Next()
	}
}

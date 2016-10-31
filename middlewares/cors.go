package middlewares

import (
	// 3rd party packages
	"github.com/gin-gonic/gin"
)

/**
 *	CORS middleware, sets required access and response headers.
 *	Aborts and returns void and sends "204 No Body" status on OPTIONS preflight request.
 *
 *	@return gin.HandlerFunc
 */
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, WWW-Authenticate, Accept, Origin, Cache-Control, X-Requested-With")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

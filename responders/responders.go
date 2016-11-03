package responders

import (
	// Native packages
	"strconv"

	// 3rd party packages
	"github.com/gin-gonic/gin"
)

type Response map[string]interface{}

/**
 *	Sends JSON response.
 *
 *	@param ctx *gin.Context
 *	@param httpStatus int
 *	@param response Response
 *
 *	@return void
 */
func ResponseObject(ctx *gin.Context, httpStatus int, response Response) {
	ctx.JSON(httpStatus, response)
}

/**
 *	Sends JSON text response, property is either "message" or "error" depending on HTTP status code.
 *
 *	@param ctx *gin.Context
 *	@param httpStatus int
 *	@param responseText string
 *
 *	@return void
 */
func ResponseText(ctx *gin.Context, httpStatus int, responseText string) {
	httpStatusSegment := httpStatus / 100
	if httpStatusSegment == 4 || httpStatusSegment == 5 {
		ctx.JSON(httpStatus, Response{
			"error": responseText,
		})
	} else {
		ctx.JSON(httpStatus, Response{
			"message": responseText,
		})
	}
}

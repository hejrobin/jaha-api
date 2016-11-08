package responders

import (
	// 3rd party packages
	"github.com/gin-gonic/gin"
)

type responseTextPrototype struct{}

func (prototype responseTextPrototype) Success(ctx *gin.Context, responseText string) {
	ResponseText(ctx, 200, responseText)
}

func (prototype responseTextPrototype) Created(ctx *gin.Context, responseText string) {
	ResponseText(ctx, 201, responseText)
}

func (prototype responseTextPrototype) BadRequest(ctx *gin.Context, responseText string) {
	ResponseText(ctx, 400, responseText)
}

func (prototype responseTextPrototype) Unauthorized(ctx *gin.Context, responseText string) {
	ResponseText(ctx, 401, responseText)
}

func (prototype responseTextPrototype) Forbidden(ctx *gin.Context, responseText string) {
	ResponseText(ctx, 403, responseText)
}

func (prototype responseTextPrototype) NotFound(ctx *gin.Context, responseText string) {
	ResponseText(ctx, 404, responseText)
}

func (prototype responseTextPrototype) Conflict(ctx *gin.Context, responseText string) {
	ResponseText(ctx, 409, responseText)
}

func (prototype responseTextPrototype) ServerError(ctx *gin.Context, responseText string) {
	ResponseText(ctx, 500, responseText)
}

func (prototype responseTextPrototype) NotImplemented(ctx *gin.Context, responseText string) {
	ResponseText(ctx, 501, responseText)
}

func Text() responseTextPrototype {
	var textResponder responseTextPrototype
	return textResponder
}

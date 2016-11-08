package responders

import (
	// 3rd party packages
	"github.com/gin-gonic/gin"
)

type responseJsonPrototype struct{}

func (prototype responseJsonPrototype) Success(ctx *gin.Context, responseObject interface{}) {
	ResponseObject(ctx, 200, responseObject)
}

func (prototype responseJsonPrototype) Created(ctx *gin.Context, responseObject interface{}) {
	ResponseObject(ctx, 201, responseObject)
}

func (prototype responseJsonPrototype) BadRequest(ctx *gin.Context, responseObject interface{}) {
	ResponseObject(ctx, 400, responseObject)
}

func (prototype responseJsonPrototype) Unauthorized(ctx *gin.Context, responseObject interface{}) {
	ResponseObject(ctx, 401, responseObject)
}

func (prototype responseJsonPrototype) Forbidden(ctx *gin.Context, responseObject interface{}) {
	ResponseObject(ctx, 403, responseObject)
}

func (prototype responseJsonPrototype) NotFound(ctx *gin.Context, responseObject interface{}) {
	ResponseObject(ctx, 404, responseObject)
}

func (prototype responseJsonPrototype) Conflict(ctx *gin.Context, responseObject interface{}) {
	ResponseObject(ctx, 409, responseObject)
}

func (prototype responseJsonPrototype) ServerError(ctx *gin.Context, responseObject interface{}) {
	ResponseObject(ctx, 500, responseObject)
}

func (prototype responseJsonPrototype) NotImplemented(ctx *gin.Context, responseObject interface{}) {
	ResponseObject(ctx, 501, responseObject)
}

func Json() responseJsonPrototype {
	var jsonResponser responseJsonPrototype
	return jsonResponser
}

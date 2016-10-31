package routers

import (
	// 3rd party packages
	"github.com/gin-gonic/gin"

	// Local packages
	"jaha-api/controllers"
	"jaha-api/middlewares"
)

/**
 *	Returns default router instance.
 *
 *	@return *gin.Engine
 */
func GetDefaultRouter() *gin.Engine {
	var router *gin.Engine

	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.Cors())

	attachDefaultRoutes(router)

	return router
}

/**
 *	Attaches route definitons to default router.
 *
 *	@param router gin.Engine - Router instance.
 *
 *	@return void
 */
func attachDefaultRoutes(router *gin.Engine) {
	router.NoRoute(controllers.DefaultController().MissingRoute)
	router.GET("", controllers.DefaultController().Index)
}

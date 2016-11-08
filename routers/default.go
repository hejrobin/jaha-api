package routers

import (
	// 3rd party packages
	"github.com/gin-gonic/gin"

	// Local packages
	"jaha-api/controllers"
	"jaha-api/env"
	"jaha-api/middlewares"
)

/**
 *	Returns default router instance.
 *
 *	@return *gin.Engine
 */
func GetDefaultRouter(sessionManager gin.HandlerFunc) *gin.Engine {
	var router *gin.Engine

	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessionManager)
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

	v1 := router.Group("v1")
	{

		Auth := middlewares.AuthMiddleware()

		v1.POST("auth", Auth.LoginHandler)

		if env.IsProductionMode() {
			v1.Use(Auth.MiddlewareFunc())
		}

		v1.Use(middlewares.Constraints())

		auth := v1.Group("auth")
		{
			auth.GET("refresh", Auth.RefreshHandler)
		}

		user := v1.Group("users")
		{
			user.GET("", controllers.UsersController().Index)
			user.POST("", controllers.UsersController().Create)

			user.GET(":uuid", controllers.UsersController().Show)
			user.PATCH(":uuid", controllers.UsersController().Update)
			user.DELETE(":uuid", controllers.UsersController().Destroy)
			user.PUT(":uuid", controllers.UsersController().Restore)
		}

		category := v1.Group("categories")
		{
			category.GET("", controllers.CategoriesController().Index)
			category.POST("", controllers.CategoriesController().Create)

			category.GET(":uuid", controllers.CategoriesController().Show)
			category.PATCH(":uuid", controllers.CategoriesController().Update)
			category.DELETE(":uuid", controllers.CategoriesController().Destroy)
			category.PUT(":uuid", controllers.CategoriesController().Restore)
		}

		statement := v1.Group("statements")
		{
			statement.GET("", controllers.StatementsController().Index)
			statement.POST("", controllers.StatementsController().Create)

			statement.GET(":uuid", controllers.StatementsController().Show)
			statement.PATCH(":uuid", controllers.StatementsController().Update)
			statement.DELETE(":uuid", controllers.StatementsController().Destroy)
			statement.PUT(":uuid", controllers.StatementsController().Restore)
		}
	}

}

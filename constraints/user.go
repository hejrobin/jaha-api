package constraints

import (
	// 3rd party packages
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	// Local packages
	"jaha-api/db"
	"jaha-api/models"
)

/**
 *	Register function for user specific callbacks.
 *	@NOTE This function *must* be called manually in router.
 *
 *	@return void
 */
func UserConstraints() {

	AddConstraint("GET", "/v1/users", func(ctx *gin.Context) bool {
		var user models.User

		session := sessions.Default(ctx)
		userId := session.Get("userId")

		if userId != "" {
			db.GetConnection().First(&user, userId)
			if user.ID != 0 && user.Role != models.USER_ROLE_ADMIN {
				return false
			}
		}

		return true
	})

}

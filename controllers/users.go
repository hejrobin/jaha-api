package controllers

import (
	// Native packages
	"fmt"
	"strconv"

	// 3rd party packages
	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
	"gopkg.in/guregu/null.v3"

	// Local packages
	"jaha-api/db"
	"jaha-api/models"
	"jaha-api/responders"
	"jaha-api/utils"
)

type usersPrototype struct{}

/**
 *	Lists published resources.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (usersPrototype) Index(ctx *gin.Context) {
	var users models.Users
	var collection models.Collection
	var collectionCount int
	var queryError error

	params := ctx.Request.URL.Query()
	paramPage, _ := strconv.Atoi(utils.Pick(params.Get("page"), "1"))
	paramOrderBy := utils.Pick(params.Get("orderBy"), "createdAt:asc")

	dbc := db.GetConnection()

	dbc.Model(&models.User{}).Count(&collectionCount)

	collection = models.Collection{}
	collection.SetLimit(COLLECTION_DEFAULT_LIMIT)
	collection.Grab(nil, 1, collectionCount)
	collection.SetPointer(paramPage)

	// Set orderBy conditions
	// @TODO Move MapOrderByConditions to Collection.
	orderByConditions := utils.MapOrderByConditions(paramOrderBy)
	query := utils.FilterOrderByConditions(dbc, models.User{}, orderByConditions)

	queryError = query.Limit(collection.Limit).Offset(collection.GetOffset()).Find(&users).Error

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	collection.Grab(users, paramPage, collectionCount)

	responders.Json().Success(ctx, collection)
	return
}

/**
 *	Retrieves published resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (usersPrototype) Show(ctx *gin.Context) {
	var user models.User
	var queryError error

	paramId := ctx.Param("uuid")
	queryError = db.GetConnection().Where("`uuid` = ?", paramId).First(&user).Error

	if user.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("User#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	responders.Json().Success(ctx, user)
	return
}

/**
 *	Creates a new resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (usersPrototype) Create(ctx *gin.Context) {
	var user models.User
	var existing models.User
	var createError error

	ctx.BindJSON(&user)

	dbc := db.GetConnection()

	dbc.Unscoped().Where("`email` = ?", user.Email).First(&existing)

	if existing.ID != 0 {
		responders.Text().BadRequest(ctx, fmt.Sprintf("Could not create resource, User#%s already exists.", existing.UUID))
		return
	}

	mergo.Merge(&user, models.User{
		Role:    1,
		UUID:    utils.RandomString(8),
		AuthKey: utils.RandomString(16),
	})

	if user.Password != "" {
		user.Password = utils.PasswordCreate(user.Password)
	}

	if !user.Valid() {
		responders.Json().BadRequest(ctx, responders.Response{
			"error":  "Resource validation failed, see issues",
			"issues": user.GetErrors(),
		})
		return
	}

	createError = dbc.Create(&user).Error

	if createError != nil {
		responders.Text().ServerError(ctx, "Could not create resource, unknown error.")
		return
	}

	responders.Json().Success(ctx, user)
	return
}

/**
 *	Updates existing resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (usersPrototype) Update(ctx *gin.Context) {
	var user models.User
	var payload models.UserPayload
	var queryError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Unscoped().Where("`uuid` = ?", paramId).First(&user).Error

	if user.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("User#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	ctx.BindJSON(&payload)

	if payload == (models.UserPayload{}) {
		responders.Text().BadRequest(ctx, "Payload cannot be empty or malformed.")
		return
	}

	if (payload.Password != "" && payload.PasswordConfirm != "") && payload.Password != payload.PasswordConfirm {
		responders.Json().BadRequest(ctx, responders.Response{
			"error": "Passwords must match.",
		})
		return
	} else {
		payload.Password = utils.PasswordCreate(payload.Password)
		payload.PasswordConfirm = payload.Password
	}

	validationError, validationErrors := utils.Validate(payload)

	if validationError != nil {
		responders.Json().BadRequest(ctx, responders.Response{
			"error":  "Resource validation failed, see issues",
			"issues": validationErrors,
		})
		return
	}

	updateError := dbc.Model(&user).Unscoped().Updates(payload).Error

	if updateError != nil {
		responders.Text().ServerError(ctx, fmt.Sprintf("Could not update User#%s.", paramId))
		return
	}

	responders.Json().Success(ctx, user)
	return
}

/**
 *	Destroys existing resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (usersPrototype) Destroy(ctx *gin.Context) {
	var user models.User
	var queryError error
	var destroyError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Unscoped().Where("`uuid` = ?", paramId).First(&user).Error

	if user.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("User#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	destroyError = dbc.Delete(&user).Error

	if destroyError != nil {
		responders.Text().ServerError(ctx, fmt.Sprintf("Could not destroy resource User#%s.", paramId))
		return
	}

	responders.NoContent(ctx)
	return
}

/**
 *	Restores soft deleted resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (usersPrototype) Restore(ctx *gin.Context) {
	var user models.User
	var queryError error
	var restoreError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Unscoped().Where("`uuid` = ?", paramId).First(&user).Error

	if user.DeletedAt == (null.Time{}) {
		responders.Text().Conflict(ctx, fmt.Sprintf("User#%s already restored.", paramId))
		return
	}

	if user.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("User#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	restoreError = dbc.Model(&user).Unscoped().Update("deleted_at", "NULL").Error

	if restoreError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	responders.Json().Success(ctx, user)
	return
}

func UsersController() usersPrototype {
	var controllerInstance usersPrototype
	return controllerInstance
}

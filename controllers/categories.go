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

type categoriesPrototype struct{}

/**
 *	Lists published resources.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (categoriesPrototype) Index(ctx *gin.Context) {
	var categories models.Categories
	var collection models.Collection
	var collectionCount int
	var queryError error

	params := ctx.Request.URL.Query()
	paramPage, _ := strconv.Atoi(utils.Pick(params.Get("page"), "1"))
	paramOrderBy := utils.Pick(params.Get("orderBy"), "createdAt:asc")

	dbc := db.GetConnection()

	dbc.Model(&models.Category{}).Count(&collectionCount)

	collection = models.Collection{}
	collection.SetLimit(COLLECTION_DEFAULT_LIMIT)
	collection.Grab(nil, 1, collectionCount)
	collection.SetPointer(paramPage)

	// Set orderBy conditions
	// @TODO Move MapOrderByConditions to Collection.
	orderByConditions := utils.MapOrderByConditions(paramOrderBy)
	query := utils.FilterOrderByConditions(dbc, models.Category{}, orderByConditions)

	queryError = query.Limit(collection.Limit).Offset(collection.GetOffset()).Find(&categories).Error

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	collection.Grab(categories, paramPage, collectionCount)

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
func (categoriesPrototype) Show(ctx *gin.Context) {
	var category models.Category
	var queryError error

	paramId := ctx.Param("uuid")
	queryError = db.GetConnection().Where("`uuid` = ?", paramId).First(&category).Error

	if category.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("Category#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	responders.Json().Success(ctx, category)
	return
}

/**
 *	Creates a new resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (categoriesPrototype) Create(ctx *gin.Context) {
	var category models.Category
	var existing models.Category
	var createError error

	ctx.BindJSON(&category)

	dbc := db.GetConnection()

	dbc.Unscoped().Where("`name` = ?", category.Name).First(&existing)

	if existing.ID != 0 {
		responders.Text().BadRequest(ctx, fmt.Sprintf("Could not create resource, Category#%s already exists.", existing.UUID))
		return
	}

	mergo.Merge(&category, models.Category{
		UUID: utils.RandomString(8),
	})

	if !category.Valid() {
		responders.Json().BadRequest(ctx, responders.Response{
			"error":  "Resource validation failed, see issues",
			"issues": category.GetErrors(),
		})
		return
	}

	createError = dbc.Create(&category).Error

	if createError != nil {
		responders.Text().ServerError(ctx, createError.Error())
		return
	}

	responders.Json().Success(ctx, category)
	return
}

/**
 *	Updates existing resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (categoriesPrototype) Update(ctx *gin.Context) {
	var category models.Category
	var payload models.CategoryPayload
	var queryError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Unscoped().Where("`uuid` = ?", paramId).First(&category).Error

	if category.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("Category#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	ctx.BindJSON(&payload)

	if payload == (models.CategoryPayload{}) {
		responders.Text().BadRequest(ctx, "Payload cannot be empty or malformed.")
		return
	}

	validationError, validationErrors := utils.Validate(payload)

	if validationError != nil {
		responders.Json().BadRequest(ctx, responders.Response{
			"error":  "Resource validation failed, see issues",
			"issues": validationErrors,
		})
		return
	}

	updateError := dbc.Model(&category).Unscoped().Updates(payload).Error

	if updateError != nil {
		responders.Text().ServerError(ctx, fmt.Sprintf("Could not update Category#%s.", paramId))
		return
	}

	responders.Json().Success(ctx, category)
	return
}

/**
 *	Destroys existing resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (categoriesPrototype) Destroy(ctx *gin.Context) {
	var category models.Category
	var queryError error
	var destroyError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Unscoped().Where("`uuid` = ?", paramId).First(&category).Error

	if category.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("Category#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	destroyError = dbc.Delete(&category).Error

	if destroyError != nil {
		responders.Text().ServerError(ctx, fmt.Sprintf("Could not destroy resource Category#%s.", paramId))
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
func (categoriesPrototype) Restore(ctx *gin.Context) {
	var category models.Category
	var queryError error
	var restoreError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Unscoped().Where("`uuid` = ?", paramId).First(&category).Error

	if category.DeletedAt == (null.Time{}) {
		responders.Text().Conflict(ctx, fmt.Sprintf("Category#%s already restored.", paramId))
		return
	}

	if category.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("Category#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	restoreError = dbc.Model(&category).Unscoped().Update("deleted_at", "NULL").Error

	if restoreError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	responders.Json().Success(ctx, category)
	return
}

func CategoriesController() categoriesPrototype {
	var controllerInstance categoriesPrototype
	return controllerInstance
}

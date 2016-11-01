package controllers

import (
	// Native packages
	"fmt"
	"strconv"

	// 3rd party packages
	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"

	// Local packages
	"jaha-api/db"
	"jaha-api/models"
	"jaha-api/scopes"
	"jaha-api/utils"
)

const COLLECTION_DEFAULT_LIMIT = 25

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

	// Get total count
	dbc.Model(&models.Categories{}).Scopes(scopes.Category().Published).Count(&collectionCount)

	// Set collection
	collection = models.Collection{
		Limit: COLLECTION_DEFAULT_LIMIT,
		Count: collectionCount,
	}

	collection.SetPointer(paramPage)
	collection.SetPageCount(collection.GetPageCount())

	query := dbc.Scopes(scopes.Category().Published)

	// Set orderBy conditions
	orderByConditions := utils.MapOrderByConditions(paramOrderBy)
	query = utils.FilterOrderByConditions(query, models.Category{}, orderByConditions)

	queryError = query.Limit(collection.Limit).Offset(collection.GetOffset()).Find(&categories).Error

	if queryError != nil {
		ctx.JSON(500, gin.H{
			"error": queryError,
		})
		return
	}

	ctx.JSON(200, models.CategoryCollection{
		Collection: categories,
		Meta:       collection,
	})
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

	dbc := db.GetConnection()
	queryError = dbc.Scopes(scopes.Category().Published).Where("`uuid` = ?", paramId).First(&category).Error

	if category.ID == 0 {
		ctx.JSON(404, gin.H{
			"error": fmt.Sprintf("Category#%s not found.", paramId),
		})
		return
	}

	if queryError != nil {
		ctx.JSON(500, gin.H{
			"error": queryError,
		})
		return
	}

	ctx.JSON(200, category)
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

	dbc.Unscoped().Where("`uuid` = ?", category.UUID).Or("`name` = ?", category.Name).Or("`slug` = ?", category.Slug).First(&existing)

	if existing.ID != 0 {
		ctx.JSON(400, gin.H{
			"error": fmt.Sprintf("Could not create resource, Category#%s already exists.", existing.UUID),
		})
		return
	}

	mergo.Merge(&category, models.Category{
		UUID: utils.RandomString(8),
	})

	validationError, validationErrors := models.Validate(category)

	if validationError != nil {
		ctx.JSON(400, gin.H{
			"error":            "Resource validation failed, see validationErrors",
			"validationErrors": validationErrors,
		})
		return
	}

	createError = dbc.Create(&category).Error

	if createError != nil {
		ctx.JSON(500, gin.H{
			"error": "Could not create resource, unknown error.",
		})
		return
	}

	ctx.JSON(201, category)
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
		ctx.JSON(404, gin.H{
			"error": fmt.Sprintf("Category#%s not found.", paramId),
		})
		return
	}

	if queryError != nil {
		ctx.JSON(500, gin.H{
			"error": queryError,
		})
		return
	}

	ctx.BindJSON(&payload)

	if payload == (models.CategoryPayload{}) {
		ctx.JSON(400, gin.H{
			"error": "Payload cannot be empty or malformed.",
		})
		return
	}

	validationError, validationErrors := models.Validate(payload)

	if validationError != nil {
		ctx.JSON(400, gin.H{
			"error":            "Resource validation failed, see validationErrors",
			"validationErrors": validationErrors,
		})
		return
	}

	updateError := dbc.Model(&category).Unscoped().Updates(payload).Error

	if updateError != nil {
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("Could not update Category#%s.", paramId),
		})
		return
	}

	ctx.JSON(200, category)
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
		ctx.JSON(404, gin.H{
			"error": fmt.Sprintf("Category#%s not found.", paramId),
		})
		return
	}

	if queryError != nil {
		ctx.JSON(500, gin.H{
			"error": queryError,
		})
		return
	}

	destroyError = dbc.Delete(&category).Error

	if destroyError != nil {
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("Could not destroy resource Category#%s.", paramId),
		})
		return
	}

	ctx.AbortWithStatus(204)
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
	queryError = dbc.Unscoped().Scopes(scopes.Category().Deleted).Where("`uuid` = ?", paramId).First(&category).Error

	if category.ID == 0 {
		ctx.JSON(404, gin.H{
			"error": fmt.Sprintf("Category#%s not found.", paramId),
		})
		return
	}

	if queryError != nil {
		ctx.JSON(500, gin.H{
			"error": queryError,
		})
		return
	}

	restoreError = dbc.Model(&category).Unscoped().Update("deleted_at", "NULL").Error

	if restoreError != nil {
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("Could not restore resource Category#%s.", paramId),
		})
		return
	}

	ctx.JSON(200, category)
	return
}

/**
 *	Returns instanciated "controller".
 *	@NOTE Classes aren't present in Go, return a struct with field methods instead.
 *
 *	@return categoriesPrototype
 */
func CategoriesController() categoriesPrototype {
	var controllerInstance categoriesPrototype
	return controllerInstance
}

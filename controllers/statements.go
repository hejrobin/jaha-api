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
	"jaha-api/scopes"
	"jaha-api/utils"
)

type statementsPrototype struct{}

/**
 *	Lists published resources.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (statementsPrototype) Index(ctx *gin.Context) {
	var statements models.Statements
	var collection models.Collection
	var collectionCount int
	var queryError error

	params := ctx.Request.URL.Query()
	paramPage, _ := strconv.Atoi(utils.Pick(params.Get("page"), "1"))
	paramOrderBy := utils.Pick(params.Get("orderBy"), "createdAt:asc")
	paramScope := params.Get("scope")

	dbc := db.GetConnection()

	// Get total count
	dbc.Model(&models.Statement{}).Count(&collectionCount)

	// Set collection
	collection = models.Collection{
		Limit: COLLECTION_DEFAULT_LIMIT,
		Count: collectionCount,
	}

	collection.SetPointer(paramPage)
	collection.SetPageCount(collection.GetPageCount())

	query := dbc.Preload("Category")
	validScope := false

	if paramScope != "" {
		switch paramScope {
		case "random":
		case "randomPick":
			validScope = true

			collection.SetPointer(1)
			collection.SetPageCount(1)

			randomLimit := collection.Limit

			if paramScope == "randomPick" {
				randomLimit = 1
			}

			queryError = query.Scopes(scopes.Statement().Random).Limit(randomLimit).Find(&statements).Error

			collection.SetCount(len(statements))
			break
		}
	}

	if !validScope {
		// Set orderBy conditions
		orderByConditions := utils.MapOrderByConditions(paramOrderBy)
		query = utils.FilterOrderByConditions(query, models.Statement{}, orderByConditions)

		queryError = query.Limit(collection.Limit).Offset(collection.GetOffset()).Find(&statements).Error
	}

	if queryError != nil {
		ctx.JSON(500, gin.H{
			"error": queryError,
		})
		return
	}

	ctx.JSON(200, models.StatementCollection{
		Collection: statements,
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
func (statementsPrototype) Show(ctx *gin.Context) {
	var statement models.Statement
	var queryError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Preload("Category").Where("`uuid` = ?", paramId).First(&statement).Error

	if statement.ID == 0 {
		ctx.JSON(404, gin.H{
			"error": fmt.Sprintf("Statement#%s not found.", paramId),
		})
		return
	}

	if queryError != nil {
		ctx.JSON(500, gin.H{
			"error": queryError,
		})
		return
	}

	ctx.JSON(200, statement)
	return
}

/**
 *	Creates a new resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (statementsPrototype) Create(ctx *gin.Context) {
	var payload models.StatementPayload
	var statement models.Statement
	var existing models.Statement
	var createError error

	ctx.BindJSON(&payload)

	dbc := db.GetConnection()

	dbc.Unscoped().Where("`body` = ?", payload.Body).First(&existing)

	if existing.ID != 0 {
		ctx.JSON(400, gin.H{
			"error": fmt.Sprintf("Could not create resource, Statement#%s already exists.", existing.UUID),
		})
		return
	}

	if payload.Category == "" {
		ctx.JSON(400, gin.H{
			"error": "Could not create resource, Category#<UUID> missing.",
		})
		return
	}

	category := models.Category{}
	categoryError := dbc.Model(&models.Category{}).Where("`uuid` = ?", payload.Category).First(&category).Error

	if categoryError != nil {
		ctx.JSON(404, gin.H{
			"error": fmt.Sprintf("Category#%s not found.", payload.Category),
		})
		return
	}

	mergo.Merge(&statement, models.Statement{
		UUID:     utils.RandomString(8),
		Body:     payload.Body,
		Category: category,
	})

	validationError, validationErrors := models.Validate(statement)

	if validationError != nil {
		ctx.JSON(400, gin.H{
			"error":            "Resource validation failed, see validationErrors",
			"validationErrors": validationErrors,
		})
		return
	}

	createError = dbc.Create(&statement).Error

	if createError != nil {
		ctx.JSON(500, gin.H{
			"error": "Could not create resource, unknown error.",
		})
		return
	}

	ctx.JSON(201, statement)
	return
}

/**
 *	Updates existing resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (statementsPrototype) Update(ctx *gin.Context) {
	var statement models.Statement
	var payload models.StatementPayload
	var queryError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Preload("Category").Unscoped().Where("`uuid` = ?", paramId).First(&statement).Error

	if statement.ID == 0 {
		ctx.JSON(404, gin.H{
			"error": fmt.Sprintf("Statement#%s not found.", paramId),
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

	if payload == (models.StatementPayload{}) {
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

	updateError := dbc.Model(&statement).Unscoped().Updates(payload).Error

	if updateError != nil {
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("Could not update Statement#%s.", paramId),
		})
		return
	}

	ctx.JSON(200, statement)
	return
}

/**
 *	Destroys existing resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (statementsPrototype) Destroy(ctx *gin.Context) {
	var statement models.Statement
	var queryError error
	var destroyError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Unscoped().Where("`uuid` = ?", paramId).First(&statement).Error

	if statement.ID == 0 {
		ctx.JSON(404, gin.H{
			"error": fmt.Sprintf("Statement#%s not found.", paramId),
		})
		return
	}

	if queryError != nil {
		ctx.JSON(500, gin.H{
			"error": queryError,
		})
		return
	}

	destroyError = dbc.Delete(&statement).Error

	if destroyError != nil {
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("Could not destroy resource Statement#%s.", paramId),
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
func (statementsPrototype) Restore(ctx *gin.Context) {
	var statement models.Statement
	var queryError error
	var restoreError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Preload("Category").Unscoped().Where("`uuid` = ?", paramId).First(&statement).Error

	if statement.DeletedAt == (null.Time{}) {
		ctx.JSON(400, gin.H{
			"error": fmt.Sprintf("Statement#%s already restored.", paramId),
		})
		return
	}

	if statement.ID == 0 {
		ctx.JSON(404, gin.H{
			"error": fmt.Sprintf("Statement#%s not found.", paramId),
		})
		return
	}

	if queryError != nil {
		ctx.JSON(500, gin.H{
			"error": queryError,
		})
		return
	}

	restoreError = dbc.Model(&statement).Unscoped().Update("deleted_at", "NULL").Error

	if restoreError != nil {
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("Could not restore resource Statement#%s.", paramId),
		})
		return
	}

	ctx.JSON(200, statement)
	return
}

/**
 *	Returns instanciated "controller".
 *	@NOTE Classes aren't present in Go, return a struct with field methods instead.
 *
 *	@return statementsPrototype
 */
func StatementsController() statementsPrototype {
	var controllerInstance statementsPrototype
	return controllerInstance
}

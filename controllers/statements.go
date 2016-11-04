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
	"jaha-api/scopes"
	"jaha-api/utils"
)

type statementsProtoype struct{}

/**
 *	Lists published resources.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (statementsProtoype) Index(ctx *gin.Context) {
	var statements models.Statements
	var collection models.Collection
	var collectionCount int
	var queryError error

	params := ctx.Request.URL.Query()
	paramPage, _ := strconv.Atoi(utils.Pick(params.Get("page"), "1"))
	paramOrderBy := utils.Pick(params.Get("orderBy"), "createdAt:asc")
	paramScope := params.Get("scope")

	dbc := db.GetConnection()

	dbc.Model(&models.Statement{}).Count(&collectionCount)

	collection = models.Collection{}
	collection.SetLimit(COLLECTION_DEFAULT_LIMIT)
	collection.Grab(nil, 1, collectionCount)
	collection.SetPointer(paramPage)

	query := dbc.Preload("Category")
	validScope := false

	if paramScope != "" {
		switch paramScope {
		case "random":
		case "randomPick":
			validScope = true

			randomLimit := collection.Limit

			if paramScope == "randomPick" {
				randomLimit = 1
			}

			queryError = query.Scopes(scopes.Statement().Random).Limit(randomLimit).Find(&statements).Error

			collection.Grab(statements, 1, len(statements))
			break
		}
	}

	if !validScope {
		// Set orderBy conditions
		// @TODO Move MapOrderByConditions to Collection.
		orderByConditions := utils.MapOrderByConditions(paramOrderBy)
		query = utils.FilterOrderByConditions(query, models.Statement{}, orderByConditions)

		queryError = query.Limit(collection.Limit).Offset(collection.GetOffset()).Find(&statements).Error

		collection.SetRecords(statements)
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

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
func (statementsProtoype) Show(ctx *gin.Context) {
	var statement models.Statement
	var queryError error

	paramId := ctx.Param("uuid")
	queryError = db.GetConnection().Preload("Category").Where("`uuid` = ?", paramId).First(&statement).Error

	if statement.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("Statement#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	responders.Json().Success(ctx, statement)
	return
}

/**
 *	Creates a new resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (statementsProtoype) Create(ctx *gin.Context) {
	var payload models.StatementPayload
	var statement models.Statement
	var existing models.Statement
	var createError error

	ctx.BindJSON(&payload)

	dbc := db.GetConnection()

	dbc.Unscoped().Where("`body` = ?", payload.Body).First(&existing)

	if existing.ID != 0 {
		responders.Text().BadRequest(ctx, fmt.Sprintf("Could not create resource, Statement#%s already exists.", existing.UUID))
		return
	}

	if payload.Category == "" {
		responders.Text().BadRequest(ctx, "Could not create resource, Category#<UUID> missing.")
		return
	}

	category := models.Category{}
	categoryError := dbc.Model(&models.Category{}).Where("`uuid` = ?", payload.Category).First(&category).Error

	if categoryError != nil {
		responders.Text().NotFound(ctx, fmt.Sprintf("Category#%s not found.", payload.Category))
		return
	}

	mergo.Merge(&statement, models.Statement{
		UUID:     utils.RandomString(8),
		Body:     payload.Body,
		Category: category,
	})

	createError = dbc.Create(&statement).Error

	if createError != nil {
		responders.Text().ServerError(ctx, "Could not create resource, unknown error.")
		return
	}

	responders.Json().Success(ctx, statement)
	return
}

/**
 *	Updates existing resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (statementsProtoype) Update(ctx *gin.Context) {
	var statement models.Statement
	var payload models.StatementPayload
	var queryError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Unscoped().Where("`uuid` = ?", paramId).First(&statement).Error

	if statement.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("Statement#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	ctx.BindJSON(&payload)

	if payload == (models.StatementPayload{}) {
		responders.Text().BadRequest(ctx, "Payload cannot be empty or malformed.")
		return
	}

	validationError, validationErrors := utils.Validate(payload)

	if validationError != nil {
		responders.Json().BadRequest(ctx, responders.Response{
			"error":  "Resource validation failed, see issues",
			"issues": validationErrors,
		})
	}

	updateError := dbc.Model(&statement).Unscoped().Updates(payload).Error

	if updateError != nil {
		responders.Text().ServerError(ctx, fmt.Sprintf("Could not update Statement#%s.", paramId))
		return
	}

	responders.Json().Success(ctx, statement)
	return
}

/**
 *	Destroys existing resource.
 *
 *	@param ctx gin.Context - Gin context pointer.
 *
 *	@return void
 */
func (statementsProtoype) Destroy(ctx *gin.Context) {
	var statement models.Statement
	var queryError error
	var destroyError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Unscoped().Where("`uuid` = ?", paramId).First(&statement).Error

	if statement.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("Statement#%s not found.", paramId))
		return
	}

	if queryError != nil {
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	destroyError = dbc.Delete(&statement).Error

	if destroyError != nil {
		responders.Text().ServerError(ctx, fmt.Sprintf("Could not destroy resource Statement#%s.", paramId))
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
func (statementsProtoype) Restore(ctx *gin.Context) {
	var statement models.Statement
	var queryError error
	var restoreError error

	paramId := ctx.Param("uuid")

	dbc := db.GetConnection()
	queryError = dbc.Unscoped().Where("`uuid` = ?", paramId).First(&statement).Error

	if statement.DeletedAt == (null.Time{}) {
		responders.Text().Conflict(ctx, fmt.Sprintf("Statement#%s already restored.", paramId))
		return
	}

	if statement.ID == 0 {
		responders.Text().NotFound(ctx, fmt.Sprintf("Statement#%s not found.", paramId))
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
		responders.Text().ServerError(ctx, queryError.Error())
		return
	}

	responders.Json().Success(ctx, statement)
	return
}

func StatementsController() statementsProtoype {
	var controllerInstance statementsProtoype
	return controllerInstance
}

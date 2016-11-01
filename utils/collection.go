package utils

import (
	// Native packages
	"fmt"
	"strings"

	// 3rd party packages
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"github.com/serenize/snaker"
)

/**
 *	Generates a map of conditions using this format "dbColumn:<asc|desc>,..."
 *
 *	@example
 *		foo:asc,bar:desc >> [foo asc, bar desc]
 *
 *	@param orderByConditions string - Conditions formatted as "dbColumn:<asc|desc>,..."
 *
 *	@return map[string]string
 */
func MapOrderByConditions(orderByConditions string) map[string]string {
	var conditions []string

	orderBy := make(map[string]string)

	for _, orderConditions := range strings.Split(orderByConditions, ",") {
		conditions = append(conditions, orderConditions)
	}

	for _, condition := range conditions {
		conditionParts := strings.Split(condition, ":")

		if len(conditionParts) > 1 {
			orderByDirection := string(conditionParts[1])
			if orderByDirection == "asc" || orderByDirection == "desc" {
				orderBy[string(conditionParts[0])] = strings.ToUpper(orderByDirection)
			}
		}
	}

	return orderBy
}

/**
 *	Filters out conditions from {@see utils.MapOrderByConditions} not present in specified model.
 *
 *	@param query *gorm.DB - Query object from GORM.
 *	@param model interface{} - Assumes model struct.
 *	@param orderByConditions map[string]string - Map with conditions generated from {@see utils.MapOrderByConditions}
 *
 *	@return *gorm.DB
 */
func FilterOrderByConditions(query *gorm.DB, model interface{}, orderByConditions map[string]string) *gorm.DB {
	modelMap := structs.Map(model)

	for conditionCol, conditionDir := range orderByConditions {
		columnFieldName := strings.Title(conditionCol)
		_, hasField := modelMap[columnFieldName]
		if hasField == true || conditionCol == "id" {
			query = query.Order(fmt.Sprintf("%s %s", snaker.CamelToSnake(conditionCol), conditionDir))
		}
	}

	return query
}

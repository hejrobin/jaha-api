package env

import (
	// Native packages
	"os"

	// Local packages
	"jaha-api/utils"
)

const ENV_DEVELOPMENT = "development"
const ENV_PRODUCTION = "production"

/**
 *	Returns MODE environment variable, defaults to ENV_DEVELOPMENT.
 *
 *	@return string
 */
func GetMode() string {
	return utils.Pick(os.Getenv("MODE"), ENV_DEVELOPMENT)
}

/**
 *	Returns PORT environment variable, defaults to 4000.
 *
 *	@return string
 */
func GetPort() string {
	return utils.Pick(os.Getenv("PORT"), "4000")
}

/**
 *	Returns true if MODE is set to ENV_PRODUCTION.
 *
 *	@return bool
 */
func IsProductionMode() bool {
	return GetMode() == ENV_PRODUCTION
}

/**
 *	Returns true if MODE is set to ENV_DEVELOPMENT.
 *
 *	@return bool
 */
func IsDevelopmentMode() bool {
	return GetMode() == ENV_DEVELOPMENT
}

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
 *	Returns database source name string.
 *
 *	@return string
 */
func GetDatabaseSourceName() string {
	return os.Getenv("DSN")
}

/**
 *	Returns database driver name string.
 *
 *	@return string
 */
func GetDatabaseDriverName() string {
	return utils.Pick(os.Getenv("DB"), "mysql")
}

/**
 *	Returns realm key used for auth realm.
 *
 *	@return string
 */
func GetRealmKey() string {
	defaultRealmKey := "!NoT--R4nD0m__R34Lm_K3Y@"
	realmKey := utils.Pick(os.Getenv("REALM"), defaultRealmKey)

	if IsProductionMode() && realmKey == defaultRealmKey {
		panic("Realm key is not set!")
	}

	return realmKey
}

/**
 *	Returns session store key.
 *
 *	@return string
 */
func GetSessionKey() string {
	defaultSessionKey := "!NoT--R4nD0m__s3ßioN_K3Y@"
	sessionKey := utils.Pick(os.Getenv("SESSION"), defaultSessionKey)

	if IsProductionMode() && sessionKey == defaultSessionKey {
		panic("Session key is not set!")
	}

	return sessionKey
}

/**
 *	Returns app name, used for auth realm etc.
 *
 *	@return string
 */
func GetAppName() string {
	return utils.Pick(os.Getenv("APP_NAME"), "jaha-api-app")
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

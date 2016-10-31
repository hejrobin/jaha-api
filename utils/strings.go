package utils

/**
 *	Returns expectedString if not empty, or fallbacks to defaultString value.
 *	@NOTE Helper function for ternery string checks in Go.
 *
 *	@param expectedString string - String to expect not to be empty.
 *	@param defaultString string - String to return if expectedString is empty.
 *
 *	@return string
 */
func Pick(expectedString string, defaultString string) string {
	if expectedString == "" {
		return defaultString
	}
	return expectedString
}

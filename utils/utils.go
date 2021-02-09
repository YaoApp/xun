package utils

// GetIF if the condition is true return the given value, else return the default value
func GetIF(condition bool, value interface{}, defaultValue interface{}) interface{} {
	if condition {
		return value
	}
	return defaultValue
}

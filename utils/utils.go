package utils

// GetIF if the condition is true return the given value, else return the default value
func GetIF(condition bool, value interface{}, defaultValue interface{}) interface{} {
	if condition {
		return value
	}
	return defaultValue
}

// PanicIF if the given value is not nil then panic, else do nothing
func PanicIF(v interface{}) {
	if v != nil {
		panic(v)
	}
}

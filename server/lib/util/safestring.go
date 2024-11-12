package util

func SafeString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
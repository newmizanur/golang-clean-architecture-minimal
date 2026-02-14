package converter

func stringValueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

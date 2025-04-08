package utils

func PString(s string) *string {
	return &s
}

func PBool(b bool) *bool {
	return &b
}

func PFloat64(f float64) *float64 {
	return &f
}

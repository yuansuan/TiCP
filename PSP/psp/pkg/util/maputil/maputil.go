package maputil

func MergeMaps[T string | int64, R any](maps ...map[T]R) map[T]R {
	mergedMap := make(map[T]R, 0)
	if len(maps) > 0 {
		for _, m := range maps {
			for key, value := range m {
				mergedMap[key] = value
			}
		}
	}

	return mergedMap
}

func EqualMaps(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v := range m1 {
		value, ok := m2[k]
		if !ok || v != value {
			return false
		}
	}

	return true
}

func ConvertSliceToMap[T string | int64](slice []T) map[T]bool {
	m := make(map[T]bool)
	for _, item := range slice {
		m[item] = true
	}
	return m
}

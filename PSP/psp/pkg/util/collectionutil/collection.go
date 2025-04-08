package collectionutil

// RemoveDuplicates 删除重复元素
func RemoveDuplicates(slice []string) []string {

	seen := make(map[string]bool)
	var result []string

	for _, element := range slice {
		if !seen[element] {
			seen[element] = true
			result = append(result, element)
		}
	}

	return result
}

func RemoveStrings(slice []string, removeSlice []string) []string {
	removeMap := make(map[string]bool)
	for _, str := range removeSlice {
		removeMap[str] = true
	}

	newSlice := []string{}
	for _, str := range slice {
		if !removeMap[str] {
			newSlice = append(newSlice, str)
		}
	}

	return newSlice
}

func RemoveString(slice []string, removeStr string) []string {
	newSlice := []string{}
	for _, item := range slice {
		if item != removeStr {
			newSlice = append(newSlice, item)
		}
	}
	return newSlice
}

func MergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

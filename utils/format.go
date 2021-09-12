package utils

// RemoveDuplicateStringValues ["aaa", "vvv", "aaa"] => ["aaa", "vvv"]
func RemoveDuplicateStringValues(source []string) []string {
	elements := make(map[string]bool)
	list := []string{}
	for _, str := range source {
		if _, exists := elements[str]; !exists {
			elements[str] = true
			list = append(list, str)
		}
	}
	return list
}

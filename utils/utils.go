package utils

// FindByteIndex e.g. 'a' in 'bab' index is 1
func FindByteIndex(find []byte, from []byte) int {
	for i, bt := range from {
		if bt == find[0] {
			isMatch := true
			for fi, fbt := range find {
				if fbt != from[i+fi] {
					isMatch = false
					break
				}
			}
			if isMatch {
				return i
			}
		}
	}
	return -1
}

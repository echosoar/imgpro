package utils

import "strconv"

// BytesToInt []byte to int
func BytesToInt(bytes []byte) int {
	byteLen := len(bytes)
	res := 0
	for index, bt := range bytes {
		res += int(bt) << ((byteLen - index - 1) * 8)
	}
	return res
}

// Uint32ToInt uint32 to int
func Uint32ToInt(num uint32) int {
	return int(num >> 8)
}

// IntToString int to string
func IntToString(num int) string {
	return strconv.Itoa(num)
}

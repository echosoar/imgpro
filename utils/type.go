package utils

import (
	"strings"
	"strconv"
)

// BytesToInt []byte to int
func BytesToInt(bytes []byte, isLow bool) int {
	byteLen := len(bytes)
	res := 0
	for index, bt := range bytes {
		if (isLow) {
			res += int(bt) << (index * 8);
		} else {
			res += int(bt) << ((byteLen - index - 1) * 8)
		}
	}
	return res
}

func ByteToRational64uString(bytes []byte, isLow bool) string {
	parts := make([]string, 0);
	for i:=0;i <len(bytes);i+=8 {
		first := BytesToInt(bytes[i:i+4], isLow);
		second := BytesToInt(bytes[i+4:i+8], isLow);
		if first == 0 {
			parts = append(parts, "0");
		} else if second == 1 {
			parts = append(parts, IntToString(first));
		} else {
			parts = append(parts, IntToString(first) + "/" + IntToString(second));
		}
	}
	return strings.Join(parts, " ");
}

// Uint32ToInt uint32 to int
func Uint32ToInt(num uint32) int {
	return int(num >> 8)
}

// IntToString int to string
func IntToString(num int) string {
	return strconv.Itoa(num)
}

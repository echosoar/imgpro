package method

import (
	"math"
	"strconv"
)

func CheckAllIsNum(list []int, num int) bool {
	for _, item := range list {
		if item != num {
			return false
		}
	}
	return true
}

func XOR(list []int, list2 []int) []int {
	length := len(list)
	if len(list2) > length {
		length = len(list2)
	}
	result := make([]int, length)
	for i := 0; i < length; i++ {
		a := 0
		b := 0
		if i < len(list) {
			a = list[i]
		}
		if i < len(list2) {
			b = list2[i]
		}
		result[i] = a ^ b
	}
	return result
}

func BinaryToInt(list []int) int {
	num := 0
	for i, value := range list {
		if value != 0 {
			num += int(math.Pow(2, float64(len(list)-1-i)))
		}

	}
	return num
}

func IntToBinary(num int, minLength int) []int {
	list := []int{}
	for num > 0 {
		list = append(list, num%2)
		num /= 2
	}
	for i := len(list); i < minLength; i++ {
		list = append(list, 0)
	}
	return ReverseArray(list)
}

// 123 => [1,2,3]
func IntToIntList(num int) []int {
	list := []int{}
	for num > 10 {
		list = append(list, num%10)
		num /= 10
	}
	list = append(list, num)
	return ReverseArray(list)
}

func IntListToString(numList []int) string {
	str := ""
	for _, item := range numList {
		str += strconv.Itoa(item)
	}
	return str
}

func Permutations(list []int, m int) [][]int {
	permutations := [][]int{}
	for index, value := range list {
		others := []int{}
		for index2, value2 := range list {
			if index2 <= index {
				continue
			}
			others = append(others, value2)
		}
		values := [][]int{}
		if m > 0 {
			child := Permutations(others, m-1)
			for _, childItems := range child {
				values = append(values, ConcatArray([]int{value}, childItems))
			}
		} else {
			values = [][]int{
				{value},
			}
		}
		permutations = append(permutations, values...)
	}
	return permutations
}

func AllPermutations(list []int) [][]int {
	allPermutations := [][]int{}
	for index := range list {
		allPermutations = append(allPermutations, Permutations(list, index)...)
	}
	return allPermutations
}

package xstr

import (
	"strconv"
	"strings"
	"unicode"
)

func ContainsLetter(password string) bool {
	for _, char := range password {
		if unicode.IsLetter(char) {
			return true
		}
	}
	return false
}

func ContainsUpperAndLower(password string) bool {
	containsUpper := false
	containsLower := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			containsUpper = true
		} else if unicode.IsLower(char) {
			containsLower = true
		}

		if containsUpper && containsLower {
			return true
		}
	}

	return false
}
func ContainsNumber(password string) bool {
	for _, char := range password {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func ContainsSpecialChar(password, specialChars string) bool {
	//specialChars := "~`!@#$%^&*()-_+={}[]|\\:;\"'<>,.?/"
	for _, char := range password {
		if strings.ContainsRune(specialChars, char) {
			return true
		}
	}
	return false
}

func ContainsChinese(str string) bool {
	for _, char := range str {
		if char >= '\u4e00' && char <= '\u9fff' {
			return true
		}
	}
	return false
}
func String2Int64(strArr []string) []int64 {
	res := make([]int64, len(strArr))

	for index, val := range strArr {
		value, _ := strconv.Atoi(val)
		res[index] = int64(value)
	}
	return res
}
func String2UInt64(strArr []string) []uint64 {
	res := make([]uint64, len(strArr))

	for index, val := range strArr {
		value, _ := strconv.Atoi(val)
		res[index] = uint64(value)
	}

	return res
}

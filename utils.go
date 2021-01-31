package main

import (
	"strconv"
)

// ExtractInt extracts an integer value from a given string slice at an index
func ExtractInt(arr []string, i int) (int, bool) {
	if len(arr) <= i {
		return 0, false
	}
	x, err := strconv.Atoi(arr[i])
	if err != nil {
		return 0, false
	}
	return x, true
}

// SafeIndexStr returns the string at index i from an array after doing out-of-bounds check.
func ExtractStr(arr []string, i int) (string, bool) {
	if len(arr) <= i {
		return "", false
	}
	return arr[i], true
}

// // StrConvert converts a slice of any type to a slice of string using fmt.Sprint
// func StrConvert(elems ...interface{}) []string {
// 	strSlc := make([]string, 0)
// 	for _, x := range elems {
// 		strSlc = append(strSlc, fmt.Sprint(x))
// 	}
// 	return strSlc
// }

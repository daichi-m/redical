package main

import "strconv"

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
func SafeIndexStr(arr []string, i int) (string, bool) {
	if len(arr) <= i {
		return "", false
	}
	return arr[i], true
}

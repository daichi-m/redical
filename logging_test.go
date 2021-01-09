package main

import (
	"math/rand"
	"testing"
)

func BenchmarkLogSafeSlice(b *testing.B) {

	lrgStr := make([]string, 0, 10000)
	lrgInt := make([]int, 0, 10000)
	lrgMap := make(map[string]int, 10000)

	for i := 0; i < 10000; i++ {
		x := rand.Int()
		y := randString(15)
		lrgInt = append(lrgInt, x)
		lrgStr = append(lrgStr, y)
		lrgMap[y] = x
	}

	tests := []struct {
		name string
		args interface{}
	}{
		{name: "string slice", args: lrgStr},
		{name: "int slice", args: lrgInt},
		{name: "single string", args: "Hello, world"},
		{name: "single int", args: 15},
		{name: "large map", args: lrgMap},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				LogSafeSlice(tt.args)
			}
		})
	}
}

func randString(len int) string {
	if len == 0 {
		return ""
	}
	ascii := rune(97) + rune(rand.Intn(26))
	return string(ascii) + randString(len-1)
}

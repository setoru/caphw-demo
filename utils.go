package main

import (
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	var result []byte
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		result = append(result, chars[rand.Intn(len(chars))])
	}
	return string(result)
}

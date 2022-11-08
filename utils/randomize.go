package utils

import "math/rand"

var payloads = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Randomize(size int) string {
	str := make([]rune, size)

	for i := range str {
		str[i] = payloads[rand.Intn(len(payloads))]
	}

	return string(str)
}

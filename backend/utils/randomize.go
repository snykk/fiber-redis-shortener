package utils

import (
	"crypto/rand"
	"log"
)

var payloads = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("Unable to generate random bytes")
	}
	return randomBytes
}

func SecureRandomString(length int) string {
	// Compute bitMask
	payloadLength := len(payloads)
	if payloadLength == 0 || payloadLength > 256 {
		panic("payloads length must be greater than 0 and less than or equal to 256")
	}
	var bitLength byte
	var bitMask byte
	for bits := payloadLength - 1; bits != 0; {
		bits = bits >> 1
		bitLength++
	}
	bitMask = 1<<bitLength - 1

	// Compute bufferSize
	bufferSize := length + length/3

	// Create random string
	result := make([]byte, length)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			// Random byte buffer is empty, get a new one
			randomBytes = generateRandomBytes(bufferSize)
		}
		// Mask bytes to get an index into the character slice
		if idx := int(randomBytes[j%length] & bitMask); idx < payloadLength {
			result[i] = payloads[idx]
			i++
		}
	}

	return string(result)
}

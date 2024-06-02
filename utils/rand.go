package utils

import (
	cyrptoRand "crypto/rand"
	"math/rand"
)

var (
	Alpha   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	Numeric = []rune("0123456789")
	Symbols = []rune("!@#$%^&*()_+{}|:<>?")
)

func RandFromSampleSpace(length int, sampleSpace []rune) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = sampleSpace[rand.Intn(len(sampleSpace))]
	}
	return string(b)
}

func RandInt(length int) int {
	return rand.Intn(length)
}

func RandBytes(arr []byte) (int, error) {
	return cyrptoRand.Read(arr)
}

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnñopqrstuvwxyzABCDEFGHIJKLMNÑOPQRSTUVWXYZ0123456789"

func main() {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	fmt.Println(string(b))
}

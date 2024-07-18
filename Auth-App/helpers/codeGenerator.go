package helpers

import (
	"math/rand"
	"time"
)

var randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomNumber() int {

	return randomGenerator.Intn(1000)
}

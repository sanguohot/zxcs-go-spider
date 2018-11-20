package random

import (
	"math/rand"
	"time"
)

func GenerateRandomIntInRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max - min) + min
	return randNum
}
func GenerateRandomFloat() float64 {
	rand.Seed(time.Now().Unix())
	return rand.Float64()
}
func SleepWithDefaultRange()  {
	time.Sleep(time.Millisecond * time.Duration(GenerateRandomIntInRange(1000, 20000)))
}
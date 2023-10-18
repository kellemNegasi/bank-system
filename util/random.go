package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nameLength = 8

var currencies = []string{"USD", "GBP", "EU"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(min, max int) int {
	return min + rand.Intn(max)
}

func RandString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		ch := letters[RandInt(1, 26)]
		sb.WriteByte(ch)
	}
	return sb.String()
}

func RandCurrency() string {
	return currencies[RandInt(0, 2)]
}

func RandomOwner() string {
	return RandString(nameLength)
}

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomMoney(min, max float64) string {
	return fmt.Sprintf("%.2f", RandomFloat(min, max))
}

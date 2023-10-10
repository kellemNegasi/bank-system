package util

import (
	"math/rand"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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

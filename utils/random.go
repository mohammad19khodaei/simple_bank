package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

const chars = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		c := chars[rand.Intn(len(chars))]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return randomString(6)
}

func RandomMoney() int64 {
	return RandomInt(1000, 5000)
}

func RandomCurrency() string {
	currencies := GetValidCurrencies()
	return currencies[rand.Intn(len(currencies))]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", randomString(6))
}

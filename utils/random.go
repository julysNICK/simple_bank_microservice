package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomINt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := (alphabet[rand.Intn(k)])
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomINt(20, 1000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR}
	return currencies[rand.Intn(len(currencies))]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}

func RandomHash() string {
	return RandomString(10)
}

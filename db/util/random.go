package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var rnd *rand.Rand

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(src)

}

func RandomInt(min, max int64) int64 {
	return min + rnd.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	i := len(alphabet)
	for j := 0; j < n; j++ {
		c := alphabet[rnd.Intn(i)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(10)
}

func RandomAmount() int {
	return int(RandomInt(0, 1000))
}

func RandomCurrency() string {
	currency := []string{USD, EUR, CAD}
	return currency[rnd.Intn(len(currency))]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

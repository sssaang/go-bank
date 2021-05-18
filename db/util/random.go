package util

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

func RandomOwner() string {
	names := []string{"Julie", "James", "Tom", "Amy", "Claire", "Alex", "Gene", "Joseph", "Tim"}
	return names[rand.Intn(len(names))]
}

func RandomMoney() int64 {
	return RandomInt(0, 10000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, KRW}
	return currencies[rand.Intn(len(currencies))]
}

func RandomEmail(name string) string {
	return fmt.Sprintf("%s@email.com", name)
}
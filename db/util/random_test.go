package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	min := int64(10)
	max := int64(10000)
	for i := 0; i < 100; i++ {
		randInt := RandomInt(min, max)
		require.GreaterOrEqual(t, randInt, min)
		require.LessOrEqual(t, randInt, max)
	}
}

func TestRandomOwner(t *testing.T) {
	names := []string{"Julie", "James", "Tom", "Amy", "Claire", "Alex", "Gene", "Joseph", "Tim"}
	for i := 0; i < 100; i++ {
		randName := RandomOwner()
		require.Contains(t, names, randName)
	}
}

func TestRandomEmail(t *testing.T) {
	for i := 0; i< 100; i++ {
		name := RandomOwner()
		randEmail := RandomEmail(name)
		require.Equal(t, fmt.Sprintf("%s@email.com", name), randEmail)
	}
}
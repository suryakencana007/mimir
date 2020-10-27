package mimir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateChar(t *testing.T) {
	l := 13
	n := 0
	for n < 64 {
		s := GenerateChar(l)
		assert.Equal(t, l, len(s))
		n++
	}
}

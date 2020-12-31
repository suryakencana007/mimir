package mimir

import (
	"math/rand"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
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

func TestULID(t *testing.T) {
	t.Parallel()

	t.Run("ULID", func(t *testing.T) {
		tm := time.Unix(1000000, 0)
		entropy := ulid.Monotonic(rand.New(rand.NewSource(tm.UnixNano())), 0)

		want := ulid.MustNew(ulid.Timestamp(tm), entropy)
		got := ULID()
		if got != want {
			t.Errorf("\ngot  %#v\nwant %#v", got, want)
		}

		assert.Equal(t, got, want)
		assert.Equal(t, got.Entropy(), want.Entropy())
	})
}

func TestUUID(t *testing.T) {
	t.Parallel()

	t.Run("UUID", func(t *testing.T) {
		got := UUID()
		assert.Equal(t, 36, len(got))
	})
}
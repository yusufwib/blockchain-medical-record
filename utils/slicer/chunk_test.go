package slicer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	type itemTest struct {
		ID    uint64
		Value string
	}

	data := []itemTest{
		{
			ID:    1,
			Value: "satu",
		},
		{
			ID:    2,
			Value: "dua",
		},
		{
			ID:    3,
			Value: "tiga",
		},
		{
			ID:    4,
			Value: "empat",
		},
	}

	want := [][]itemTest{
		{
			{
				ID:    1,
				Value: "satu",
			},
			{
				ID:    2,
				Value: "dua",
			},
			{
				ID:    3,
				Value: "tiga",
			},
		},
		{
			{
				ID:    4,
				Value: "empat",
			},
		},
	}

	got := Chunk(data, 3)

	assert.Equal(t, want, got)

}

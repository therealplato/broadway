package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortVars(t *testing.T) {
	foo1 := varKV{
		k: "food",
		v: "clamato",
	}

	foo2 := varKV{
		k: "fool",
		v: "motley",
	}

	foo3 := varKV{
		k: "goal",
		v: "goooooooooooooooooooooooal",
	}
	v1 := varSlice{foo1, foo2, foo3}
	v2 := varSlice{foo2, foo3, foo1}

	testcases := []struct {
		scenario string
		in       varSlice
		out      varSlice
	}{
		{
			scenario: "it sorts in order vars",
			in:       v1,
			out:      v1,
		},
		{
			scenario: "it sorts out of order vars",
			in:       v2,
			out:      v1,
		},
	}

	for _, tc := range testcases {
		out := tc.in
		sortVars(out)
		assert.Equal(t, tc.out, out, tc.scenario)
	}

}

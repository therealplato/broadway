package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortVars(t *testing.T) {
	foo1 := Var{
		k: "food",
		v: "clamato",
	}

	foo2 := Var{
		k: "fool",
		v: "motley",
	}

	foo3 := Var{
		k: "goal",
		v: "goooooooooooooooooooooooal",
	}
	v1 := Vars{foo1, foo2, foo3}
	v2 := Vars{foo2, foo3, foo1}

	testcases := []struct {
		scenario string
		in       Vars
		out      Vars
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

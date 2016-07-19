package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortVars(t *testing.T) {
	foo := Var{
		k: "foo",
		v: "bar",
	}

	goo := Var{
		k: "goo",
		v: "car",
	}

	hoo := Var{
		k: "hoo",
		v: "dar",
	}
	v1 := Vars{foo, goo, hoo}
	v2 := Vars{goo, hoo, foo}

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

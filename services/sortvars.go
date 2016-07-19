package services

import (
	"sort"
	"strings"
)

// Implement sort.Interface for Vars:

// Vars is a slice of instance variables
type Vars []Var

// Var is a single instance variable with key k, value v
type Var struct {
	k string
	v string
}

// Less returns true if vv[i].k is alphabetically before vv[j].k
func (vv Vars) Less(i, j int) bool {
	xx := []Var(vv)
	if len(xx[i].k) == 0 {
		return true
	}
	if len(xx[j].k) == 0 {
		return false
	}
	return strings.ToLower(xx[i].k)[0] < strings.ToLower(xx[j].k)[0]
}

// Len returns the length of vv
func (vv Vars) Len() int {
	xx := []Var(vv)
	return len(xx)
}

// Swap swaps two indices of vv
func (vv Vars) Swap(i, j int) {
	xx := []Var(vv)
	// *vv[i], *vv[j] = *vv[j], *vv[i]
	xx[i], xx[j] = xx[j], xx[i]
}

// sortVars takes a slice of 2-ary arrays
// each array is [k, v]
// returned slice will be sorted alphabetically by k
func sortVars(vv Vars) {
	sort.Sort(vv)
}

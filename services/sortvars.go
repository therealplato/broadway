package services

import (
	"sort"
	"strings"
)

// Implement sort.Interface for Vars:

// varSlice is a slice of instance variables
type varSlice []varKV

// varKV is a single instance variable with key k, value v
type varKV struct {
	k string
	v string
}

// Less returns true if vv[i].k is alphabetically before vv[j].k
func (vv varSlice) Less(i, j int) bool {
	xx := []varKV(vv)
	if len(xx[i].k) == 0 {
		return true
	}
	if len(xx[j].k) == 0 {
		return false
	}
	return strings.ToLower(xx[i].k) < strings.ToLower(xx[j].k)
}

// Len returns the length of vv
func (vv varSlice) Len() int {
	xx := []varKV(vv)
	return len(xx)
}

// Swap swaps two indices of vv
func (vv varSlice) Swap(i, j int) {
	xx := []varKV(vv)
	// *vv[i], *vv[j] = *vv[j], *vv[i]
	xx[i], xx[j] = xx[j], xx[i]
}

// sortVars takes a slice of 2-ary arrays
// each array is [k, v]
// returned slice will be sorted alphabetically by k
func sortVars(vv varSlice) {
	sort.Sort(vv)
}

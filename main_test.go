package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testresult := m.Run()
	os.Exit(testresult)
}

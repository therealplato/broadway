package main

import (
	"github.com/namely/broadway/fixtures"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fixtures.SetupTestFixtures()
	testresult := m.Run()
	fixtures.TeardownTestFixtures()
	os.Exit(testresult)
}

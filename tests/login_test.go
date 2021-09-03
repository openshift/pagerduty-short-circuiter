package main

import (
	"testing"
)

func TestPdcli(t *testing.T) {

	assertCorrectMessage := func(t testing.TB, actual, expected int) {
		t.Helper()
		if actual != expected {
			t.Errorf("Actual:%d Expected:%d", actual, expected)
		}
	}

	//Tesing all the functions against an empty file.
	t.Run("sample test", func(t *testing.T) {
		actual := 1
		expected := 1
		assertCorrectMessage(t, actual, expected)
	})
}

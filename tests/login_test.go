package main

import (
	"testing"
    "github.com/openshift/pagerduty-short-circuiter/cmd/pdcli/login"
)

func TestWc(t *testing.T) {

	assertCorrectMessage := func(t testing.TB, actual, expected int) {
		t.Helper()
		if actual != expected {
			t.Errorf("Actual:%d Expected:%d", actual, expected)
		}
	}

    //Tesing all the functions against an empty file.
	t.Run("login test", func(t *testing.T) {
		actual := 
		expected := "logged in"
		assertCorrectMessage(t, actual, expected)
	})}

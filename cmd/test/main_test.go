package main_test

import (
	"api_chat/api/tests"
	"testing"
)

func TestApiChat(t *testing.T) {
	tk, err := tests.NewTestsKeys(t)
	if err != nil {
		t.Errorf("error init tests keys: %v", err)
	} else {
		tk.TestDatabase()
		tk.TestRepos()
	}
}

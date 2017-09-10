package gini

import (
	"testing"
)

func TestGini(t *testing.T) {
	config, _ := Parse("test.ini")
	value := config.Get("mysql", "username")

	if value != "mysql" {
		t.Errorf("Failed: expected %s, got %s", "mysql", value)
	}
}

func TestGlobal(t *testing.T) {
	config, _ := Parse("test.ini")
	value := config.Get("", "foo")

	if value != "bar" {
		t.Errorf("Failed: expected %s, got %s", "bar", value)
	}
}

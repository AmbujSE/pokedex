package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "   mixed    CASE   TeST ",
			expected: []string{"mixed", "case", "test"},
		},
		{
			input:    "   ", // Edge case: empty input
			expected: []string{},
		},
		{
			input:    "singleword",
			expected: []string{"singleword"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("For input '%s': Expected %v, but got %v", c.input, c.expected, actual)
		}
	}
}

package main

import (
	"testing"
)


func TestCleanInput(t *testing.T) {
	cases := []struct {
	input string
	expected []string
	}{
		{
			input: "  heLLo  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "One",
			expected: []string{"one"},
		},
		{
			input: " one   two   three ",
			expected: []string{"one", "two", "three"},
		},
		{
			input: "",
			expected: nil,
		},
		{
			input: "    ",
			expected: nil,
		},
	}
	for _, c := range(cases) {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Number of words doesn't match")
		}
		
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("expected: %v -- Actual: %v", expectedWord, word)
			}
		}
	}	
}
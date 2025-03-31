package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello  world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "     hello      world      ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		output := cleanInput(c.input)
		if len(output) != len(c.expected) {
			t.Errorf("Number of words is incorrect: %d and %d", len(output), len(c.expected))
			t.FailNow()
		}
		for i := range output {
			if c.expected[i] != output[i] {
				t.Errorf("Expected %s, got %s at position %d", c.expected[i], output[i], i)
				// t.Errorf("error: words do not match")
				// t.Fail()
				t.FailNow()
			}
		}
	}

}
